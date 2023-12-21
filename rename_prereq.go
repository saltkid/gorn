package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func series_rename_prereqs(path string, s_type string, keep_ep_nums bool, starting_ep_num int, has_season_0 bool) (map[string]any, error) {
	// get prerequsite info for renaming series
	info := map[string]any {
		"path": path,
		"type": s_type,
		"keep_ep_nums": keep_ep_nums,
		"starting_ep_num": starting_ep_num,

		"ep_nums": make([]int, 0),
		"seasons": make(map[string]int),

		"has_season_0": has_season_0,
		"extras_dirs": make([]string, 0),
	}

	if s_type == "named_seasons" || s_type == "multiple_season_no_movies" || s_type == "multiple_season_with_movies" {
		seasons, extras_dirs, err := get_seasons_and_extras(path, s_type, has_season_0)
		if err != nil {
			return nil, err
		}
		info["seasons"] = seasons
		info["extras_dirs"] = extras_dirs

	} else if s_type == "single_season_no_movies" || s_type == "single_season_with_movies" {
		_, extras_dirs, err := get_seasons_and_extras(path, s_type, has_season_0)
		if err != nil {
			return nil, err
		}
		seasons := make(map[int]string)
		seasons[1] = filepath.Base(path)
		info["extras_dirs"] = extras_dirs
		info["seasons"] = seasons

	} else {
		return nil, fmt.Errorf("unknown series type: %s", s_type)
	}

	return info, nil
}

func get_seasons_and_extras (path string, s_type string, has_season_0 bool) (map[int]string, []string, error) {
	seasons := make(map[int]string)
	extras := make([]string, 0)
	
	subdirs, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	for _, subdir := range subdirs {
		if !subdir.IsDir() {
			continue
		}

		// skip subdir with same name as directory if s_type is 'single_season_with_movies'
		if s_type == "single_season_with_movies" && subdir.Name() == filepath.Base(path) {
			continue
		}

		if has_season_0 {
			if seasons[0] != "" {
				return nil, nil, fmt.Errorf("multiple seasons found in %s", path)
			}

			extras_pattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?|ova`)
			if extras_pattern.MatchString(subdir.Name()) {
				seasons[0] = subdir.Name()
				continue
			}
		}

		// skip getting season numbers if s_type is 'single_season_no_movies' or 'single_season_with_movies'
		// put whatever dir it is in extras
		if s_type == "single_season_no_movies" || s_type == "single_season_with_movies" || s_type == "multiple_season_with_movies" {
		    extras = append(extras, subdir.Name())
			continue
		}

		// get season number from subdir name
		var season_name_pattern *regexp.Regexp
		if s_type == "named_seasons" {
			season_name_pattern = regexp.MustCompile(`^(\d+)\..*$`)
		} else if s_type == "multiple_season_no_movies" || s_type == "multiple_season_with_movies" {
			season_name_pattern = regexp.MustCompile(`^(?i)season\s+(\d+).*$`)
		} else {
			return nil, nil, fmt.Errorf("unknown series type: %s; series type must be one of 'named_seasons', 'multiple_season_no_movies', 'multiple_season_with_movies'", s_type)
		}
		if season_name_pattern == nil {
			return nil, nil, fmt.Errorf("unknown series type: %s; series type must be one of 'named_seasons', 'multiple_season_no_movies', 'multiple_season_with_movies'", s_type)
		}

		season_num := season_name_pattern.FindStringSubmatch(subdir.Name())
		if season_num == nil {
			extras = append(extras, subdir.Name())
			continue
		}

		// season_num[0] is the whole string so we only need season_num[1] (first matched group)
		num, err := strconv.Atoi(season_num[1])
		if err != nil {
			return nil, nil, err
		}
		seasons[num] = subdir.Name()
	}

	return seasons, extras, nil
}