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
	keep_ep_nums    bool
	starting_ep_num int
	seasons         map[int]string
	movies          []string
	has_season_0    bool
	extras_dirs     []string
}

type MovieInfo struct {
	path        string
	movie_type  string
	movies      []string
	extras_dirs []string
}

func (info *SeriesInfo) rename() error {
	// for padding of season numbers when renaming: min 2 digits
	max_season_digits := len(strconv.Itoa(len(info.seasons)))
	if max_season_digits < 2 {
		max_season_digits = 2
	}

	for num, season := range info.seasons {
		var season_path string
		if info.series_type == "single_season_no_movies" {
			season_path = info.path
		} else if info.series_type == "single_season_with_movies" || info.series_type == "named_seasons" || info.series_type == "multiple_season_no_movies" || info.series_type == "multiple_season_with_movies" {
			season_path = info.path + "/" + season
		} else {
			return fmt.Errorf("unknown series type: %s", info.series_type)
		}

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
		if info.starting_ep_num > 0 {
			ep_num = info.starting_ep_num
		} else {
			ep_num = 1
		}

		ep_nums := make([]int, 0)
		if info.keep_ep_nums {
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
			ep_num_str := fmt.Sprintf("%0*d", max_ep_digits, ep_nums[i])
			season_num_str := fmt.Sprintf("%0*d", max_season_digits, num)

			new_name := "S" + season_num_str + "E" + ep_num_str + filepath.Ext(file)
			fmt.Println(fmt.Sprintf("%-*s", 20, file), " --> ", fmt.Sprintf("%*s", 20, new_name))
		}
		fmt.Println()
	}
	return nil
}

func (info *MovieInfo) rename() error {
	return nil
}
