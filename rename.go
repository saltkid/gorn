package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sort"
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
			var title string
			if info.series_type == "single_season_no_movies" || info.series_type == "multiple_season_no_movies" || info.series_type == "multiple_season_with_movies" {
				title = filepath.Base(info.path)
			} else if info.series_type == "single_season_with_movies" {
				title = filepath.Base(season_path)
			} else if info.series_type == "named_seasons" {
				title = filepath.Base(info.path) + " " + filepath.Base(season_path)
			}
			
			new_name := fmt.Sprintf("S%0*dE%0*d %s%s",
									max_season_digits, num, 
									max_ep_digits, ep_nums[i],
									clean_title(title), filepath.Ext(file))

			fmt.Println(fmt.Sprintf("%-*s", 20, file), " --> ", fmt.Sprintf("%*s", 20, new_name))
			fmt.Println("old", season_path+"/"+file, "new", season_path+"/"+new_name)
			_, err := os.Stat(season_path+new_name)
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
