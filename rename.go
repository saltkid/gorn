package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sort"
	"strings"
)

type Rename interface {
	rename() error
}

type SeriesInfo struct {
	path            string
	series_type     string
	seasons         map[int]string
	movies          []string
	keep_ep_nums    Option[bool]
	starting_ep_num Option[int]
	naming_scheme   Option[string]
}

type MovieInfo struct {
	path        string
	movie_type  string
	movies      map[string]string
}

func (info *SeriesInfo) rename() error {
	// for padding of season numbers when renaming: min 2 digits
	max_season_digits := len(strconv.Itoa(len(info.seasons)))
	if max_season_digits < 2 {
		max_season_digits = 2
	}

	// rename episodes
	for num, season := range info.seasons {
		is_valid_type := map[string]bool{
			"single_season_no_movies": true,
			"single_season_with_movies": true,
			"named_seasons": true,
			"multiple_season_no_movies": true,
			"multiple_season_with_movies": true,
		}
		if !is_valid_type[info.series_type]{
			return fmt.Errorf("unknown series type: %s", info.series_type)
		}

		season_path := filepath.Clean(info.path + "/" + season)

		fmt.Println("path: ", season_path)
		files, err := os.ReadDir(season_path)
		if err != nil {
			return err
		}
		media_files := make([]string, 0)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if is_media_file(file.Name()) {
				media_files = append(media_files, file.Name())
			}
		}
		sort.Sort(FilenameSort(media_files))

		max_ep_digits := len(strconv.Itoa(len(media_files)))
		if max_ep_digits < 2 {
			max_ep_digits = 2
		}
		
		var ep_num int
		sen, err := info.starting_ep_num.get()
		if err != nil {
			return err
		}

		if sen > 0 {
			ep_num = sen
		} else {
			ep_num = 1
		}

		ep_nums := make([]int, 0)
		ken, err := info.keep_ep_nums.get()
		if err != nil {
			return err
		}

		if ken {
			for _, file := range media_files {
				ep_num, err = read_episode_num(file)
				if err != nil {
					return err
				}
				
				temp_max := len(strconv.Itoa(ep_num))
				if temp_max > max_ep_digits {
					max_ep_digits = temp_max
				}

				ep_nums = append(ep_nums, ep_num)
			}
		
		} else {
			for range media_files {
				ep_nums = append(ep_nums, ep_num)
				ep_num++
			}
		}

		for i, file := range media_files {
			title := default_title(info.series_type, info.naming_scheme, info.path, season_path)
			new_name, err := generate_new_name(info.naming_scheme,
										  max_season_digits, num, 	// season_pad, season_num
										  max_ep_digits, ep_nums[i],// ep_pad, ep_num 
										  title, file) 				// title, file path
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%-*s", 20, file), " --> ", fmt.Sprintf("%*s", 20, new_name))
			fmt.Println("old", season_path+"/"+file, "new", season_path+"/"+new_name)
			_, err = os.Stat(season_path+new_name)
			if err == nil {
				fmt.Println("renaming", season_path+"/"+file, "to", season_path+"/"+new_name + " failed: file already exists")
				continue
			} else if os.IsNotExist(err) {
				err = os.Rename(season_path+"/"+file, season_path+"/"+new_name)
			} else {
				return err
			}
			if err != nil {
				return err
			}
		}
		fmt.Println()
	}

	// rename movies if needed
	if info.series_type == "single_season_with_movies" || info.series_type == "multiple_season_with_movies" {
		for _,movie := range info.movies {
			files, err := os.ReadDir(info.path + "/" + movie)
			if err != nil {
				return err
			}

			media_files := make([]string, 0)
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if is_media_file(file.Name()) {
					media_files = append(media_files, file.Name())
				}
			}

			if len(media_files) > 1 {
				return fmt.Errorf("multiple media files found in %s for a movie direcotry in %s", movie, info.path+"/"+filepath.Base(movie))
			} else if len(media_files) == 0 {
				return fmt.Errorf("no media files found in %s for a movie directory in %s", movie, info.path+"/"+filepath.Base(movie))
			}

			new_name := fmt.Sprintf("%s %s%s", filepath.Base(info.path), filepath.Base(movie), filepath.Ext(media_files[0]))
			fmt.Println(fmt.Sprintf("%-*s", 20, media_files[0]), " --> ", fmt.Sprintf("%*s", 20, new_name))
			fmt.Println("old", info.path+"/"+movie+"/"+media_files[0], "new", info.path+"/"+movie+"/"+new_name)
			err = os.Rename(info.path+"/"+movie+"/"+media_files[0], info.path+"/"+movie+"/"+new_name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (info *MovieInfo) rename() error {
	for dir, file := range info.movies {
		new_name := clean_title(dir) + filepath.Ext(file)
		old_name := file
		if info.movie_type == "movie_set" {
			old_name = dir + "/" + old_name
			new_name = dir + "/" + new_name
		}

		fmt.Println(fmt.Sprintf("%-*s", 20, old_name), " --> ", fmt.Sprintf("%*s", 20, new_name))
		fmt.Println("old", info.path+"/"+old_name, "new", info.path+"/"+new_name)
		err := os.Rename(info.path+"/"+old_name, info.path+"/"+new_name)
		if err != nil {
			return err
		}
	}
	return nil
}

func default_title(series_type string, naming_scheme Option[string], path string, season_path string) string {
	var title string
	if series_type == "single_season_no_movies" || series_type == "multiple_season_no_movies" || series_type == "multiple_season_with_movies" {
		title = filepath.Base(path)
	} else if series_type == "single_season_with_movies" {
		title = filepath.Base(season_path)
	} else if series_type == "named_seasons" {
		title = filepath.Base(path) + " " + filepath.Base(season_path)
	}
	return title
}

func generate_new_name(naming_scheme Option[string], season_pad int, season_num int, ep_pad int, ep_num int, title string, file string) (string, error) {
	var new_name string
	if naming_scheme.is_some() {
		scheme, err := naming_scheme.get()
		if err != nil {
			return "", err
		}
		// replace <season_num>
		new_name = regexp.MustCompile(`<season_num(\s*:\s*\d+)?>`).ReplaceAllStringFunc(scheme, func(match string) string {
			// <season_num: \d+>
			if strings.Contains(match, ":") {
				pad := regexp.MustCompile(`\d+`).FindString(match)
				pad_num, err := strconv.Atoi(pad)
				if err != nil {
					return match
				}
				return fmt.Sprintf("%0*d", pad_num, season_num)
			}
			// <season_num>
			return fmt.Sprintf("%0*d", season_pad, season_num)
		})
		// replace <episode_num>
		new_name = regexp.MustCompile(`<episode_num(\s*:\s*\d+)?>`).ReplaceAllStringFunc(new_name, func(match string) string {
			// <episode_num: \d+>
			if strings.Contains(match, ":") {
				pad := regexp.MustCompile(`\d+`).FindString(match)
				pad_num, err := strconv.Atoi(pad)
				if err != nil {
					return match
				}
				return fmt.Sprintf("%0*d", pad_num, ep_num)
			}
			// <episode_num>
			return fmt.Sprintf("%0*d", ep_pad, ep_num)
		})
		// replace <self: start,end> with filepath.Base(file)[start:end]
		new_name = regexp.MustCompile(`<self\s*:\s*\d+,\d+>`).ReplaceAllStringFunc(new_name, func(match string) string {
			parts := regexp.MustCompile(`\d+`).FindAllString(match, 2)
			if len(parts) != 2 {
				return match
			}
			start, err := strconv.Atoi(parts[0])
			if err != nil {
				return match
			}
			end, err := strconv.Atoi(parts[1])
			if err != nil {
				return match
			}
			fmt.Println(filepath.Base(file)[start:end])
			return filepath.Base(file)[start:end]
		})
		// replace <parent> tokens with nth parent's name
		// lol goodluck: https://regex-vis.com/?r=%3C%28parent%28-parent%29*%28%5Cs*%3A%5Cs*%28%28%5Cd%2B%5Cs*%2C%5Cs*%5Cd%2B%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%7Cp%28-%5Cd%2B%29%3F%28%5Cs*%3A%5Cs*%28%28%5Cd%2B%5Cs*%2C%5Cs*%5Cd%2B%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%29%5Cs*%3E&e=0
		new_name = regexp.MustCompile(`<(parent(-parent)*(\s*:\s*((\d+\s*,\s*\d+)|('[^']*')))?|p(-\d+)?(\s*:\s*((\d+\s*,\s*\d+)|('[^']*')))?)\s*>`).ReplaceAllStringFunc(new_name, func(match string) string {
			n, err := parent_token_to_int(match)
			if err != nil {
				return match
			}
			parent_name := nth_parent(file, n)

			if strings.Contains(match, ":") {
				parts := regexp.MustCompile(`\d+`).FindAllString(match, 2)

				if len(parts) == 2 {
					start, err := strconv.Atoi(parts[0])
					if err != nil {
						return parent_name
					}
					end, err := strconv.Atoi(parts[1])
					if err != nil {
						return parent_name
					}
					return parent_name[start:end]
				}
				regex_pattern := regexp.MustCompile(`'[^']*'`).FindString(match)
				if len(regex_pattern) > 0 {
					regex_pattern = strings.Trim(regex_pattern, "'")
					re, err := regexp.Compile(regex_pattern)
					if err != nil {
						return parent_name
					}
					
					substrings := re.FindStringSubmatch(parent_name)
					if len(substrings) > 1 {
						return substrings[1]
					} 
				}
			}
			return parent_name
		})
		// append ext
		new_name = fmt.Sprintf("%s%s", new_name, filepath.Ext(file))

	} else {
		new_name = fmt.Sprintf("S%0*dE%0*d %s%s",
							season_pad, season_num, 
							ep_pad, ep_num,
							title, filepath.Ext(file))
	}

	return new_name, nil
}