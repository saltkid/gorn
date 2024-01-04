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
	options         AdditionalOptions
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

		var media_files []string
		err := filepath.WalkDir(season_path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && is_media_file(d.Name()) {
				media_files = append(media_files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
		sort.Sort(FilenameSort(media_files))

		max_ep_digits := len(strconv.Itoa(len(media_files)))
		if max_ep_digits < 2 {
			max_ep_digits = 2
		}
		
		// if additional options are none aka user inputted var, ask for user input
		season_options := prompt_additional_options(info.options, season_path, 2)

		var ep_num, sen int
		if season_options.starting_ep_num.is_some() {
			sen, _ = season_options.starting_ep_num.get()
		} else {
			sen = 1
		}
		if sen > 0 {
			ep_num = sen
		} else {
			ep_num = 1
		}

		ep_nums := make([]int, 0)
		var ken bool
		if season_options.keep_ep_nums.is_some() {
			ken, _ = season_options.keep_ep_nums.get()
		} else {
			ken = false
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
			title := default_title(info.series_type, season_options.naming_scheme, info.path, season_path)
			new_name, err := generate_new_name(season_options.naming_scheme,// naming_scheme
											   max_season_digits, num, 		// season_pad, season_num
										  	   max_ep_digits, ep_nums[i],	// ep_pad, ep_num 
										  	   title, file)					// title, file path
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%-*s", 20, file), " --> ", fmt.Sprintf("%*s", 20, new_name))
			fmt.Println("old", file, "\nnew", new_name)
			_, err = os.Stat(new_name)
			if err == nil {
				fmt.Println("renaming", filepath.Base(file), "to", filepath.Base(new_name) + " failed: file already exists")
				continue
			} else if os.IsNotExist(err) {
				err = os.Rename(file, new_name)
				if err != nil {
					return err
				}
			} else {
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
	return clean_title(title)
}

func generate_new_name(naming_scheme Option[string], season_pad int, season_num int, ep_pad int, ep_num int, title string, abs_path string) (string, error) {
	var new_name string
	ns, _ := naming_scheme.get()
	if naming_scheme.is_some() && ns != "default" {
		scheme, err := naming_scheme.get()
		if err != nil {
			return "", err
		}
		// replace <season_num>
		new_name = regexp.MustCompile(`<season_num(\s*:\s*\d+)?>`).ReplaceAllStringFunc(scheme, func(match string) string {
			// <season_num: \d+>
			if strings.Contains(match, ":") {
				pad := regexp.MustCompile(`\d+`).FindString(match)
				pad_num, _ := strconv.Atoi(pad)
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
				pad_num, _ := strconv.Atoi(pad)
				return fmt.Sprintf("%0*d", pad_num, ep_num)
			}
			// <episode_num>
			return fmt.Sprintf("%0*d", ep_pad, ep_num)
		})
		// replace <self>
		new_name = regexp.MustCompile(`<self\s*:\s*\d+,\d+>`).ReplaceAllStringFunc(new_name, func(match string) string {
			// if error, return full base name without extension
			base_name := filepath.Base(abs_path)
			base_name = strings.ReplaceAll(base_name, filepath.Ext(base_name), "")

			parts := regexp.MustCompile(`\d+`).FindAllString(match, 2)
			if len(parts) != 2 {
				return base_name
			}
			start, err := strconv.Atoi(parts[0])
			if err != nil || start >= len(base_name) {
				return base_name
			}
			end, err := strconv.Atoi(parts[1])
			if err != nil || end+1 >= len(base_name) {
				return base_name
			}
			return base_name[start:end+1]
		})
		// replace <parent> tokens with nth parent's name
		// lol goodluck: https://regex-vis.com/?r=%3C%28parent%28-parent%29*%28%5Cs*%3A%5Cs*%28%28%5Cd%2B%5Cs*%2C%5Cs*%5Cd%2B%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%7Cp%28-%5Cd%2B%29%3F%28%5Cs*%3A%5Cs*%28%28%5Cd%2B%5Cs*%2C%5Cs*%5Cd%2B%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%29%5Cs*%3E&e=0
		new_name = regexp.MustCompile(`<(parent(-parent)*(\s*:\s*((\d+(\s*,\s*\d+)?)|('[^']*')))?|p(-\d+)?(\s*:\s*((\d+(\s*,\s*\d+)?)|('[^']*')))?)\s*>`).ReplaceAllStringFunc(new_name, func(match string) string {
			n, err := parent_token_to_int(match)
			if err != nil {
				return new_name
			}
			parent_name := nth_parent(abs_path, n)

			// <parent>
			if !strings.Contains(match, ":") {return parent_name}

			// has ':'
			// <parent: <value>>
			trimmed_match := strings.Trim(match, "<>")
			val := strings.TrimSpace(strings.SplitN(trimmed_match, ":", 2)[1])
			switch val[0] {
			// <parent: 1,2>
			case ',':
				val := strings.SplitN(val, ",", 2)
				start, err := strconv.Atoi(val[0])
				if err != nil || start >= len(parent_name) {
					return parent_name
				}
				end, err := strconv.Atoi(val[1])
				if err != nil || end+1 >= len(parent_name) {
					return parent_name
				}
				return parent_name[start:end+1]

			// <parent: '<regex_pattern>'>
			case '\'':
				regex_pattern := strings.Trim(val, "'")
				_, err := regexp.Compile(regex_pattern)
				if err != nil {
					return parent_name
				}
				sub_regexes := split_regex_by_pipe(regex_pattern)
				for _, re := range sub_regexes {
					sub_match := regexp.MustCompile(re).FindStringSubmatch(parent_name)
					if len(sub_match) > 1 {
						// found a substring match
						return sub_match[1]
					}
				}
				// did not find a substring match
				return parent_name

			// <parent: 1>
			default:
				start, err := strconv.Atoi(val)
				if err != nil || start+1 >= len(parent_name) {
					return parent_name
				}
				return parent_name[start:start+1]
			}
		})
		// append ext
		new_name = filepath.Join(filepath.Dir(abs_path), fmt.Sprintf("%s%s", new_name, filepath.Ext(abs_path)))

	} else if naming_scheme.is_none() || ns == "default"{
		new_name = fmt.Sprintf("S%0*dE%0*d %s%s",
							season_pad, season_num, 
							ep_pad, ep_num,
							title, filepath.Ext(abs_path))
		new_name = filepath.Join(filepath.Dir(abs_path), new_name)
	}

	return new_name, nil
}