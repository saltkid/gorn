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

func SeriesRenamePrereqs(path string, sType string, options AdditionalOptions) (SeriesInfo, error) {
	// Get prerequsite info for renaming series
	isValidType := map[string]bool{
		"singleSeasonNoMovies":     true,
		"singleSeasonWithMovies":   true,
		"namedSeasons":             true,
		"multipleSeasonNoMovies":   true,
		"multipleSeasonWithMovies": true,
	}
	if !isValidType[sType] {
		return SeriesInfo{}, fmt.Errorf("unknown series type: %s", sType)
	}

	// if additional options are none aka user inputted var, ask for user input
	options = PromptOptionalFlags(options, path, 1)
	info := SeriesInfo{
		path:       path,
		seriesType: sType,
		seasons:    make(map[int]string),
		movies:     make([]string, 0),
		options:    options,
	}

	s0, err := options.hasSeason0.Get()
	if err != nil {
		return SeriesInfo{}, err
	}
	seasons, movies, err := FetchSeriesContent(path, sType, s0)
	if err != nil {
		return SeriesInfo{}, err
	}
	if sType == "singleSeasonNoMovies" {
		seasons[1] = ""
	} else if sType == "singleSeasonWithMovies" {
		seasons[1] = filepath.Base(path)
	}
	info.seasons = seasons
	info.movies = movies

	return info, nil
}

// PromptOptionalFlags prompts the user for additional options.
//
// params:
//   - options AdditionalOptions: Additional options for the prompt.
//   - path string: The path of the file.
//   - level int8: The level of the prompt.
//
// levels:
//   - level 0: per series type level
//   - level 1: per series entry level
//   - level 2: per series season level
//
// return:
//   - AdditionalOptions: The additional options for the prompt.
func PromptOptionalFlags(options AdditionalOptions, path string, level int8) AdditionalOptions {
	defaultKEN := some[bool](false)
	defaultSEN := some[int](1)
	defaultS0 := some[bool](false)
	defaultNS := some[string]("default")

	var varOpt []string
	var s0Opt []string
	if level == 0 {
		varOpt = []string{"var/", ", 'var'"}
		s0Opt = varOpt
	} else if level == 1 {
		varOpt = []string{"var/", ", 'var'"}
		s0Opt = []string{"", ""}
	} else if level == 2 {
		varOpt = []string{"", ""}
		s0Opt = varOpt
	}

	// prompt user for additional options
	if options.keepEpNums.IsNone() {
		fmt.Printf("[INPUT]\nkeep episode numbers for '%s'?\ninputs: (y/n/%sdefault/exit)\n", filepath.Base(path), varOpt[0])
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := strings.ToLower(strings.TrimSpace(scanner.Text()))

				if input == "y" || input == "yes" {
					options.keepEpNums = some[bool](true)
					break
				} else if input == "n" || input == "no" {
					options.keepEpNums = some[bool](false)
					break
				} else if input == "var" && level < 2 {
					break
				} else if input == "exit" {
					return options
				} else if input == "default" {
					options.keepEpNums = defaultKEN
					break
				} else {
					fmt.Printf("[ERROR]\ninvalid input, please enter 'y', 'n'%s, 'exit', or 'default'\n", varOpt[1])
				}
			}
		}
	}
	if options.startingEpNum.IsNone() {
		fmt.Printf("[INPUT]\nstarting episode number for '%s'?\ninputs: (<int>/%sdefault/exit)\n", filepath.Base(path), varOpt[0])
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := strings.ToLower(strings.TrimSpace(scanner.Text()))

				int_input, err := strconv.Atoi(input)
				if err == nil {
					options.startingEpNum = some[int](int_input)
					break
				}
				if input == "default" {
					options.startingEpNum = defaultSEN
					break
				} else if input == "var" && level < 2 {
					break
				} else if input == "exit" {
					return options
				} else {
					fmt.Printf("[ERROR]\ninvalid input, please enter '<int>'%s, 'exit', or 'default'\n", varOpt[1])
				}
			}
		}
	}
	if options.hasSeason0.IsNone() {
		fmt.Printf("[INPUT]\nspecials/extras directory under '%s' as season 0?\ninputs: (y/n/%sdefault/exit)\n", filepath.Base(path), s0Opt[0])
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := strings.ToLower(strings.TrimSpace(scanner.Text()))

				if input == "y" || input == "yes" {
					options.hasSeason0 = some[bool](true)
					break
				} else if input == "n" || input == "no" {
					options.hasSeason0 = some[bool](false)
					break
				} else if input == "var" && level == 0 {
					break
				} else if input == "exit" {
					return options
				} else if input == "default" {
					options.hasSeason0 = defaultS0
					break
				} else {
					fmt.Printf("[ERROR]\ninvalid input, please enter 'y', 'n'%s, 'exit', or 'default'\n", s0Opt[1])
				}
			}
		}
	}
	namingScheme, _ := options.namingScheme.Get()
	if options.namingScheme.IsNone() || namingScheme != "default" {
		fmt.Printf("[INPUT]\nnaming scheme for '%s'?\ninputs: (<naming scheme>/%sdefault)\n", filepath.Base(path), varOpt[0])
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := scanner.Text()
				input = strings.TrimSpace(input)

				if strings.ToLower(input) == "var" && level < 2 {
					break
				} else if strings.ToLower(input) == "default" {
					options.namingScheme = defaultNS
					break
				} else if strings.ToLower(input) == "exit" {
					return options
				} else if err := ValidateNamingScheme(input); err == nil && input != "var" {
					options.namingScheme = some[string](input)
					break
				} else {
					fmt.Printf("[ERROR]\ninvalid input, please enter 'y', 'n'%s, 'default', 'exit', or a valid naming scheme\n", varOpt[1])
					fmt.Println("input:", input)
					if err != nil {
						fmt.Println("naming scheme error:", err)
					} else {
						fmt.Println("error: invalid input")
					}
				}
			}
		}
	}

	return options
}

func FetchSeriesContent(path string, sType string, hasSeason0 bool) (map[int]string, []string, error) {
	seasons := make(map[int]string)
	movies := make([]string, 0)

	subdirs, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	extrasPattern := regexp.MustCompile(`^(?i)(specials?|extras?|o(v|n)a)`)
	for _, subdir := range subdirs {
		if !subdir.IsDir() {
			continue
		}

		// skip subdir with same name as directory if sType is 'singleSeasonWithMovies'
		// season is assigned outside of this function
		if sType == "singleSeasonWithMovies" && subdir.Name() == filepath.Base(path) {
			continue
		}

		if hasSeason0 {
			if seasons[0] != "" {
				return nil, nil, fmt.Errorf("multiple specials/extras directories found in %s", path)
			}

			if extrasPattern.MatchString(subdir.Name()) {
				seasons[0] = subdir.Name()
				continue
			}
		}

		if sType == "singleSeasonNoMovies" {
			continue
		} else if sType == "singleSeasonWithMovies" {
			if extrasPattern.MatchString(subdir.Name()) {
				continue
			} else {
				movies = append(movies, subdir.Name())
				continue
			}
		}

		// Get season number from subdir name
		var seasonNamePattern *regexp.Regexp
		if sType == "namedSeasons" {
			seasonNamePattern = regexp.MustCompile(`^(\d+)\..*$`)
		} else if sType == "multipleSeasonNoMovies" || sType == "multipleSeasonWithMovies" {
			seasonNamePattern = regexp.MustCompile(`^(?i)season\s+(\d+).*$`)
		} else {
			return nil, nil, fmt.Errorf("unknown series type: %s; series type must be one of 'namedSeasons', 'multipleSeasonNoMovies', 'multipleSeasonWithMovies'", sType)
		}
		if seasonNamePattern == nil {
			return nil, nil, fmt.Errorf("unknown series type: %s; series type must be one of 'namedSeasons', 'multipleSeasonNoMovies', 'multipleSeasonWithMovies'", sType)
		}

		readSeasonNum := seasonNamePattern.FindStringSubmatch(subdir.Name())
		if readSeasonNum == nil {
			if sType == "multipleSeasonWithMovies" {
				if extrasPattern.MatchString(subdir.Name()) {
					continue
				} else {
					movies = append(movies, subdir.Name())
					continue
				}
			}
			continue
		}

		// readSeasonNum[0] is the whole string so we only need readSeasonNum[1] (first matched group)
		num, err := strconv.Atoi(readSeasonNum[1])
		if err != nil {
			return nil, nil, err
		}
		seasons[num] = subdir.Name()
	}

	return seasons, movies, nil
}

func MovieRenamePrereqs(path string, mType string) (MovieInfo, error) {
	info := MovieInfo{
		path:      path,
		movieType: mType,
		movies:    make(map[string]string),
	}

	subdirs, err := os.ReadDir(path)
	if err != nil {
		return MovieInfo{}, err
	}

	extrasPattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?|ova`)

	for _, subdir := range subdirs {
		if mType == "standalone" {
			if subdir.IsDir() && extrasPattern.MatchString(subdir.Name()) {
				continue
			}

			if IsMediaFile(subdir.Name()) {
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

		if mType == "movieSet" {
			if extrasPattern.MatchString(subdir.Name()) {
				continue
			}

			files, err := os.ReadDir(filepath.Join(path, subdir.Name()))
			if err != nil {
				return MovieInfo{}, err
			}

			movieCount := 0
			for _, file := range files {
				if IsMediaFile(file.Name()) {
					if movieCount > 0 {
						return MovieInfo{}, fmt.Errorf("multiple media files found in %s", path)
					}

					info.movies[subdir.Name()] = file.Name()
					movieCount++
				}
			}

			if movieCount == 0 {
				return MovieInfo{}, fmt.Errorf("no media files found in %s", path)
			}
		}
	}

	return info, nil
}
