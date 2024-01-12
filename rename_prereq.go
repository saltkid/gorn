package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// SeriesRenamePrereqs returns a SeriesInfo with information needed to rename a series entry.
//	path: The path of the series.
//	sType: The type of series.
//	options: Flags for the prompt.
func SeriesRenamePrereqs(path string, sType string, options Flags) (SeriesInfo) {
	// if flags are none aka user inputted var, ask for user input again
	options = PromptOptionalFlags(options, path, 1)
	info := SeriesInfo{
		path:       path,
		seriesType: sType,
		seasons:    make(map[int]string),
		movies:     make([]string, 0),
		options:    options,
	}

	// has season 0 will always be of some value at this point
	s0, _ := options.hasSeason0.Get()
	seasons, movies := FetchSeriesContent(path, sType, s0)
	info.seasons = seasons
	info.movies = movies

	return info
}

// PromptOptionalFlags prompts the user for additional options.
//
// params:
//   - options Flags: flags for the prompt.
//   - path string: The path of the file.
//   - level int8: The level of the prompt.
//
// levels:
//   - level 0: per series type level
//   - level 1: per series entry level
//   - level 2: per series season level
//
// return:
//   - Flags: changed/unchanged flags for the prompt.
func PromptOptionalFlags(options Flags, path string, level int8) Flags {
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
					namingScheme := strings.Trim(input, `"`)
					options.namingScheme = some[string](namingScheme)
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

// FetchSeriesContent retrieves the season directories and movie directories from the given series entry.
//	path: The path of the series entry.
//	sType: The type of series.
//	hasSeason0: Whether the series has a season 0 directory.
func FetchSeriesContent(path string, sType string, hasSeason0 bool) (map[int]string, []string) {
	seasons := make(map[int]string)
	movies := make([]string, 0)

	subdirs, err := os.ReadDir(path)
	if err != nil {
		log.Println(WARN, "error reading series entry:", err, "; skipping renaming entry:", path)
		return seasons, movies
	}

	extrasPattern := regexp.MustCompile(`^(?i)(specials?|extras?|o(v|n)a)`)
	for _, subdir := range subdirs {
		if !subdir.IsDir() {
			continue
		}

		if sType == SINGLE_SEASON_WITH_MOVIES && subdir.Name() == filepath.Base(path) {
			seasons[1] = subdir.Name()
			continue
		}

		if hasSeason0 {
			if extrasPattern.MatchString(subdir.Name()) {
				if seasons[0] != "" {
					log.Println(WARN, "multiple specials/extras directories found [", seasons[0], ",", subdir.Name(), "]; skipping renaming entry:", path )
					return make(map[int]string), make([]string, 0)
				}
				seasons[0] = subdir.Name()
				continue
			}
		}

		if sType == SINGLE_SEASON_NO_MOVIES {
			seasons[1] = ""
			continue
		} else if sType == SINGLE_SEASON_WITH_MOVIES {
			if extrasPattern.MatchString(subdir.Name()) {
				continue
			} else {
				movies = append(movies, subdir.Name())
				continue
			}
		}

		// Get season number from subdir name
		var seasonNamePattern *regexp.Regexp
		if sType == NAMED_SEASONS {
			seasonNamePattern = regexp.MustCompile(`^(\d+)\s*(?:\.|_|-).*$`)
		} else if sType == MULTIPLE_SEASON_NO_MOVIES || sType == MULTIPLE_SEASON_WITH_MOVIES {
			seasonNamePattern = regexp.MustCompile(`^(?i)season\s+(\d+).*$`)
		}

		readSeasonNum := seasonNamePattern.FindStringSubmatch(subdir.Name())
		if readSeasonNum == nil {
			if sType == MULTIPLE_SEASON_WITH_MOVIES {
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
		// first matched group will always be a number
		num, _ := strconv.Atoi(readSeasonNum[1])
		seasons[num] = subdir.Name()
	}

	return seasons, movies
}

// MovieRenamePrereqs returns a MovieInfo with information needed to rename a movie entry.
//	path: The path of the movie.
//	mType: The type of movie.
func MovieRenamePrereqs(path string, mType string) (MovieInfo) {
	info := MovieInfo{
		path:      path,
		movieType: mType,
		movies:    make(map[string]string),
	}
	defaultInfo := info

	subdirs, err := os.ReadDir(path)
	if err != nil {
		log.Println(WARN, "error reading movie entry:", err, "; skipping renaming entry:", path)
		return info
	}

	extrasPattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?|ova`)

	for _, subdir := range subdirs {
		if mType == STANDALONE {
			if subdir.IsDir() && extrasPattern.MatchString(subdir.Name()) {
				continue
			}

			if IsMediaFile(subdir.Name()) {
				if len(info.movies) == 0 {
					info.movies[filepath.Base(path)] = subdir.Name()
					continue
				} else {
					log.Println(WARN, "multiple media files found in supposedly standalone movie directory: [", info.movies[filepath.Base(path)], ",", subdir.Name(), "]; skipping renaming entry:", path )
					return defaultInfo
				}
			}
		}

		if !subdir.IsDir() {
			continue
		}

		if mType == MOVIE_SET {
			if extrasPattern.MatchString(subdir.Name()) {
				continue
			}

			files, err := os.ReadDir(filepath.Join(path, subdir.Name()))
			if err != nil {
				log.Println(WARN, "error reading entry:", filepath.Join(path, subdir.Name()), "; skipping renaming entry:", path)
				return defaultInfo
			}

			movieCount := 0
			skipEntry := false
			for _, file := range files {
				if IsMediaFile(file.Name()) {
					if movieCount > 0 {
						log.Println(WARN, "multiple media files found in", file, ": [", info.movies[subdir.Name()], ",", file.Name(), "]; skipping renaming movie:", subdir, "under movie set:", path)
						skipEntry = true
						break
					}

					info.movies[subdir.Name()] = file.Name()
					movieCount++
				}
			}
			if skipEntry {
				continue
			}

			if movieCount == 0 {
				log.Println(WARN, "no media files found in", subdir.Name(), "; skipping renaming movie entry:", subdir)
				continue
			}
		}
	}

	return info
}
