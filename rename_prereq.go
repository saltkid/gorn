package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func series_rename_prereqs(path string, s_type string, options AdditionalOptions) (SeriesInfo, error) {
	// get prerequsite info for renaming series
	is_valid_type := map[string]bool{
		"single_season_no_movies": true,
		"single_season_with_movies": true,
		"named_seasons": true,
		"multiple_season_no_movies": true,
		"multiple_season_with_movies": true,
	}
	if !is_valid_type[s_type] {
		return SeriesInfo{}, fmt.Errorf("unknown series type: %s", s_type)
	}

	// if additional options are none aka user inputted var, ask for user input
	options = prompt_additional_options(options, path)
	info := SeriesInfo{
		path: 				path,
		series_type: 		s_type,
		seasons: 			make(map[int]string),
		movies: 			make([]string, 0),
		options: 			options,
	}

	s0, err := options.has_season_0.get()
	if err != nil {
		return SeriesInfo{}, err
	}
	seasons, movies, err := fetch_series_content(path, s_type, s0)
	if err != nil {
		return SeriesInfo{}, err
	}
	if s_type == "single_season_no_movies" {
		seasons[1] = ""
	} else if s_type == "single_season_with_movies" {
		seasons[1] = filepath.Base(path)
	}
	info.seasons = seasons
	info.movies = movies

	return info, nil
}

func prompt_additional_options(options AdditionalOptions, path string) (AdditionalOptions) {
	default_ken := some[bool](false)
	default_sen := some[int](1)
	default_s0 := some[bool](false)
	default_ns := some[string]("default")

	// prompt user for additional options
	if options.keep_ep_nums.is_none() {
		fmt.Println("[INPUT]\nkeep episode numbers for", filepath.Base(path), "?", "\ninputs: (y/n/var/default/exit)")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := strings.ToLower(strings.TrimSpace(scanner.Text()))

				if input == "y" || input == "yes" {
					options.keep_ep_nums = some[bool](true)
					break
				} else if input == "n" || input == "no" {
					options.keep_ep_nums = some[bool](false)
					break
				} else if input == "var" {
					break
				} else if input == "exit" {
					return options
				} else if input == "default" {
					options.keep_ep_nums = default_ken
					break
				} else {
					fmt.Println("[ERROR]\ninvalid input, please enter 'y', 'n', 'exit', or 'default'")
				}
			}
		}
	}
	if options.starting_ep_num.is_none() {
		fmt.Println("[INPUT]\nstarting episode number for", filepath.Base(path), "?\ninputs: (<int>/var/default/exit)")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := strings.ToLower(strings.TrimSpace(scanner.Text()))

				int_input, err := strconv.Atoi(input)
				if err == nil {
					options.starting_ep_num = some[int](int_input)
					break
				}
				if input == "default" {
					options.starting_ep_num = default_sen
					break
				} else if input == "var" {
					break
				} else if input == "exit" {
					return options
				} else {
					fmt.Println("[ERROR]\ninvalid input, please enter '<int>', 'var', 'exit', or 'default'")
				}
			}
		}
	}
	if options.has_season_0.is_none() {
		fmt.Println("[INPUT]\nspecials/extras directory under", filepath.Base(path), "as season 0?", "\ninputs: (y/n/var/default/exit)")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := strings.ToLower(strings.TrimSpace(scanner.Text()))

				if input == "y" || input == "yes" {
					options.has_season_0 = some[bool](true)
					break
				} else if input == "n" || input == "no" {
					options.has_season_0 = some[bool](false)
					break
				} else if input == "var" {
					break
				} else if input == "exit" {
					return options
				} else if input == "default" {
					options.has_season_0 = default_s0
					break
				} else {
					fmt.Println("[ERROR]\ninvalid input, please enter 'y', 'n', 'exit', or 'default'")
				}
			}
		}
	}
	naming_scheme, _ := options.naming_scheme.get()
	if options.naming_scheme.is_none() || naming_scheme != "default" {
		fmt.Println("[INPUT]\nnaming scheme for", filepath.Base(path), "\ninputs: (<naming scheme>/var/default)")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := scanner.Text()
				input = strings.TrimSpace(input)

				if strings.ToLower(input) == "var" {
					break
				} else if input == "default" {
					options.naming_scheme = default_ns
				} else if input == "exit" {
					return options
				} else if err := validate_naming_scheme(input); err == nil {
					options.naming_scheme = some[string](input)
					break
				} else {
					fmt.Println("[ERROR]\ninvalid input, please enter 'y', 'n', 'default', 'exit', or a valid naming scheme\ninput:", input, "\nerror:", err)
				}
			}
		}
	}

	return options
}

func fetch_series_content(path string, s_type string, has_season_0 bool) (map[int]string, []string, error) {
	seasons := make(map[int]string)
	movies := make([]string, 0)
	
	subdirs, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	extras_pattern := regexp.MustCompile(`^(?i)(specials?|extras?|o(v|n)a)`)
	for _, subdir := range subdirs {
		if !subdir.IsDir() {
			continue
		}

		// skip subdir with same name as directory if s_type is 'single_season_with_movies'
		// season is assigned outside of this function
		if s_type == "single_season_with_movies" && subdir.Name() == filepath.Base(path) {
			continue
		}

		if has_season_0 {
			if seasons[0] != "" {
				return nil, nil, fmt.Errorf("multiple specials/extras directories found in %s", path)
			}

			if extras_pattern.MatchString(subdir.Name()) {
				seasons[0] = subdir.Name()
				continue
			}
		}

		if s_type == "single_season_no_movies" {
			continue
		} else if s_type == "single_season_with_movies"{
			if extras_pattern.MatchString(subdir.Name()) {
				continue
			} else {
				movies = append(movies, subdir.Name())
				continue
			}
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
			if s_type == "multiple_season_with_movies" {
				if extras_pattern.MatchString(subdir.Name()) {
					continue
				} else {
					movies = append(movies, subdir.Name())
					continue
				}
			}
			continue
		}

		// season_num[0] is the whole string so we only need season_num[1] (first matched group)
		num, err := strconv.Atoi(season_num[1])
		if err != nil {
			return nil, nil, err
		}
		seasons[num] = subdir.Name()
	}

	return seasons, movies, nil
}

func movie_rename_prereqs(path string, m_type string) (MovieInfo, error) {
	info := MovieInfo{
		path: 			path,
		movie_type: 	m_type,
		movies: 		make(map[string]string),
	}

	subdirs, err := os.ReadDir(path)
	if err != nil {
		return MovieInfo{}, err
	}

	extras_pattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?|ova`)

	for _, subdir := range subdirs {
		if m_type == "standalone" {
			if subdir.IsDir() && extras_pattern.MatchString(subdir.Name()) {
				continue
			}

			if is_media_file(subdir.Name()) {
				if len(info.movies) == 0 {
					info.movies[filepath.Base(path)] = subdir.Name()
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

					info.movies[subdir.Name()] = file.Name()
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