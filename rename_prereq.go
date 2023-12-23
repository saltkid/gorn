package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type SeriesInfo struct {
	path string
	series_type string
	keep_ep_nums bool
	starting_ep_num int
	seasons map[int]string
	movies []string
	has_season_0 bool
	extras_dirs []string
}

type MovieInfo struct {
	path string
	movie_type string
	movies []string
	extras_dirs []string
}

func series_rename_prereqs(path string, s_type string, keep_ep_nums bool, starting_ep_num int, has_season_0 bool) (SeriesInfo, error) {
	// get prerequsite info for renaming series
	info := SeriesInfo{
		path: 				path,
		series_type: 		s_type,
		keep_ep_nums: 		keep_ep_nums,
		starting_ep_num: 	starting_ep_num,
		seasons: 			make(map[int]string),
		movies: 			make([]string, 0),
		has_season_0: 		has_season_0,
		extras_dirs: 		make([]string, 0),
	}

	if s_type == "named_seasons" || s_type == "multiple_season_no_movies" || s_type == "multiple_season_with_movies" {
		seasons, extras_dirs, movies, err := get_series_content(path, s_type, has_season_0)
		if err != nil {
			return SeriesInfo{}, err
		}
		info.seasons = seasons
		info.movies = movies
		info.extras_dirs = extras_dirs

	} else if s_type == "single_season_no_movies" || s_type == "single_season_with_movies" {
		_, extras_dirs, movies, err := get_series_content(path, s_type, has_season_0)
		if err != nil {
			return SeriesInfo{}, err
		}
		seasons := make(map[int]string)
		seasons[1] = filepath.Base(path)
		info.seasons = seasons
		info.movies = movies
		info.extras_dirs = extras_dirs

	} else {
		return SeriesInfo{}, fmt.Errorf("unknown series type: %s", s_type)
	}

	return info, nil
}

func get_series_content (path string, s_type string, has_season_0 bool) (map[int]string, []string, []string, error) {
	seasons := make(map[int]string)
	extras := make([]string, 0)
	movies := make([]string, 0)
	
	subdirs, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, nil, err
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
				return nil, nil, nil, fmt.Errorf("multiple seasons found in %s", path)
			}

			extras_pattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?|ova`)
			if extras_pattern.MatchString(subdir.Name()) {
				seasons[0] = subdir.Name()
				continue
			}
		}

		if s_type == "single_season_no_movies" {
		    extras = append(extras, subdir.Name())
			continue
		} else if s_type == "single_season_with_movies"{
			movies = append(movies, subdir.Name())
			continue
		}

		// get season number from subdir name
		var season_name_pattern *regexp.Regexp
		if s_type == "named_seasons" {
			season_name_pattern = regexp.MustCompile(`^(\d+)\..*$`)
		} else if s_type == "multiple_season_no_movies" || s_type == "multiple_season_with_movies" {
			season_name_pattern = regexp.MustCompile(`^(?i)season\s+(\d+).*$`)
		} else {
			return nil, nil, nil, fmt.Errorf("unknown series type: %s; series type must be one of 'named_seasons', 'multiple_season_no_movies', 'multiple_season_with_movies'", s_type)
		}
		if season_name_pattern == nil {
			return nil, nil, nil, fmt.Errorf("unknown series type: %s; series type must be one of 'named_seasons', 'multiple_season_no_movies', 'multiple_season_with_movies'", s_type)
		}

		season_num := season_name_pattern.FindStringSubmatch(subdir.Name())
		if season_num == nil {
			if s_type == "multiple_season_with_movies" {
				movies = append(movies, subdir.Name())
			} else {
				extras = append(extras, subdir.Name())
			}
			continue
		}

		// season_num[0] is the whole string so we only need season_num[1] (first matched group)
		num, err := strconv.Atoi(season_num[1])
		if err != nil {
			return nil, nil, nil, err
		}
		seasons[num] = subdir.Name()
	}

	return seasons, extras, movies, nil
}

func movie_rename_prereqs (path string, m_type string) (MovieInfo, error) {
	info := MovieInfo{
		path: 			path,
		movie_type: 	m_type,
		movies: 		make([]string, 0),
		extras_dirs: 	make([]string, 0),
	}

	subdirs, err := os.ReadDir(path)
	if err != nil {
		return MovieInfo{}, err
	}

	extras_pattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?|ova`)

	for _, subdir := range subdirs {
		if m_type == "standalone" {
			if subdir.IsDir() && extras_pattern.MatchString(subdir.Name()) {
				info.extras_dirs = append(info.extras_dirs, subdir.Name())
				continue
			}

			if is_media_file(subdir.Name()) {
				if len(info.movies) == 0 {
					info.movies = append(info.movies, subdir.Name())
					continue
				} else {
					return MovieInfo{}, fmt.Errorf("multiple media files found in %s for an entry marked as a standalone movie", path)
				}
			}
		}

		if !subdir.IsDir() {
			continue
		}

		if m_type == "movie_set" {
			if extras_pattern.MatchString(subdir.Name()) {
				info.extras_dirs = append(info.extras_dirs, subdir.Name())
				continue
			}

			files, err := os.ReadDir(filepath.Join(path, subdir.Name()))
			if err != nil {
				return MovieInfo{}, err
			}

			movie_count := 0
			for _, file := range files {
				if is_media_file(file.Name()) {
					if movie_count > 0 {
						return MovieInfo{}, fmt.Errorf("multiple media files found in %s", path)
					}

					info.movies = append(info.movies, filepath.Join(subdir.Name(), file.Name()))
					movie_count++
				}
			}

			if movie_count == 0 {
				return MovieInfo{}, fmt.Errorf("no media files found in %s", path)	
			}
		}
	}

	return info, nil
}