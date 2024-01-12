package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Rename interface {
	Rename() error
}

type SeriesInfo struct {
	path       string
	seriesType string
	seasons    map[int]string
	movies     []string
	options    Flags
}

type MovieInfo struct {
	path      string
	movieType string
	movies    map[string]string
}

func (info *SeriesInfo) Rename() {
	// for padding of season numbers when renaming: min 2 digits
	maxSeasonDigits := len(strconv.Itoa(len(info.seasons)))
	if maxSeasonDigits < 2 {
		maxSeasonDigits = 2
	}

	// Rename episodes
	for num, season := range info.seasons {
		seasonPath := filepath.Clean(info.path + "/" + season)
		var mediaFiles []string
		err := filepath.WalkDir(seasonPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && IsMediaFile(d.Name()) {
				mediaFiles = append(mediaFiles, path)
			}
			return nil
		})
		if err != nil {
			log.Println(WARN, "error reading media files:", err, "; skipping renaming all episodes in:", seasonPath)
			continue
		}
		sort.Sort(FilenameSort(mediaFiles))

		maxEpDigits := len(strconv.Itoa(len(mediaFiles)))
		if maxEpDigits < 2 {
			maxEpDigits = 2
		}

		// if additional options are none aka user inputted var, ask for user input
		seasonOptions := PromptOptionalFlags(info.options, seasonPath, 2)

		var epNum, sen int
		if seasonOptions.startingEpNum.IsSome() {
			sen, _ = seasonOptions.startingEpNum.Get()
		} else {
			sen = 1
		}
		if sen > 0 {
			epNum = sen
		} else {
			epNum = 1
		}

		epNums := make([]int, 0)
		var ken bool
		if seasonOptions.keepEpNums.IsSome() {
			ken, _ = seasonOptions.keepEpNums.Get()
		} else {
			ken = false
		}

		if ken {
			for i, file := range mediaFiles {
				epNum, err = ReadEpisodeNum(file)
				if err != nil {
					log.Println(WARN, "error reading episode number from", file, ":", err, "; skipping renaming")
					// don't include this episode in renaming
					mediaFiles = append(mediaFiles[:i], mediaFiles[i+1:]...)
				}

				tempMax := len(strconv.Itoa(epNum))
				if tempMax > maxEpDigits {
					maxEpDigits = tempMax
				}

				epNums = append(epNums, epNum)
			}

		} else {
			for range mediaFiles {
				epNums = append(epNums, epNum)
				epNum++
			}
		}

		// adjust episode read episode numbers if starting episode number was specified by user
		if seasonOptions.startingEpNum.IsSome() {
			adjustVal := 0
			minEpNum := epNums[0]
			for _, val := range epNums[1:] {
				if val < minEpNum {
					minEpNum = val
				}
			}
			adjustVal = sen - minEpNum
			for i, val := range epNums {
				epNums[i] = val + adjustVal
			}
		}

		for i, file := range mediaFiles {
			title := DefaultTitle(info.seriesType, seasonOptions.namingScheme, info.path, seasonPath)
			newName := GenerateNewName(seasonOptions.namingScheme, // namingScheme
				maxSeasonDigits, num, // season_pad, season_num
				maxEpDigits, epNums[i], // ep_pad, epNum
				title, file) // title, file path
			
			// TODO: decide whether to turn this into log or not
			// fmt.Println(fmt.Sprintf("%-*s", 20, file), " --> ", fmt.Sprintf("%*s", 20, newName))
			// fmt.Println("old", file, "\nnew", newName)

			_, err = os.Stat(newName)
			if err == nil {
				log.Println(WARN, "file already exists: renaming", filepath.Base(file), "to", filepath.Base(newName), "failed; skipping renaming:", file)
				continue
			} else if os.IsNotExist(err) {
				err = os.Rename(file, newName)
				if err != nil {
					log.Println(WARN, "renaming error:", err, "; skipping renaming:", file)
					continue
				}
			} else {
				log.Println(WARN, "unexpected error when checking if file exists before renaming:", err)
				continue
			}
		}
	}

	// Rename movies if needed
	if info.seriesType == SINGLE_SEASON_WITH_MOVIES || info.seriesType == MULTIPLE_SEASON_WITH_MOVIES {
		for _, movie := range info.movies {
			files, err := os.ReadDir(info.path + "/" + movie)
			if err != nil {
				log.Println(WARN, "error reading media files under ", info.path+"/"+movie, ":", err, "; skipping renaming")
				continue
			}

			mediaFiles := make([]string, 0)
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if IsMediaFile(file.Name()) {
					mediaFiles = append(mediaFiles, file.Name())
				}
			}

			if len(mediaFiles) > 1 {
				log.Println(WARN, "multiple media files found in supposed movie directory under series:", info.path+"/"+movie, "; skipping renaming")
				continue
			} else if len(mediaFiles) == 0 {
				log.Println(WARN, "no media files found in:", info.path+"/"+movie, "; skipping renaming")
				continue
			}

			newName := fmt.Sprintf("%s %s%s", filepath.Base(info.path), filepath.Base(movie), filepath.Ext(mediaFiles[0]))
			// TODO: decide whether to turn this into log or not
			// fmt.Println(fmt.Sprintf("%-*s", 20, mediaFiles[0]), " --> ", fmt.Sprintf("%*s", 20, newName))
			// fmt.Println("old", info.path+"/"+movie+"/"+mediaFiles[0], "new", info.path+"/"+movie+"/"+newName)
			err = os.Rename(info.path+"/"+movie+"/"+mediaFiles[0], info.path+"/"+movie+"/"+newName)
			if err != nil {
				log.Println(WARN, "renaming error:", err, "; skipping renaming:", info.path+"/"+movie+"/"+mediaFiles[0])
				continue
			}
		}
	}
}

func (info *MovieInfo) Rename(wg *sync.WaitGroup) {
	for dir, file := range info.movies {
		wg.Add(1)
		go func(dir string, file string) {
			defer wg.Done()

			newName := CleanTitle(dir) + filepath.Ext(file)
			old_name := file
			if info.movieType == "movieSet" {
				old_name = dir + "/" + old_name
				newName = dir + "/" + newName
			}
	
			// TODO: decide whether to turn this into log or not
			// fmt.Println(fmt.Sprintf("%-*s", 20, old_name), " --> ", fmt.Sprintf("%*s", 20, newName))
			// fmt.Println("old", info.path+"/"+old_name, "new", info.path+"/"+newName)
			_, err := os.Stat(info.path+"/"+newName)
			if err == nil {
				log.Println(WARN, "file already exists: renaming", filepath.Base(old_name), "to", filepath.Base(newName), "failed; skipping renaming:", info.path+"/"+old_name)
				return
			} else if os.IsNotExist(err) {
				err = os.Rename(info.path+"/"+old_name, info.path+"/"+newName)
				if err != nil {
					log.Println(WARN, "renaming error:", err, "; skipping renaming:", info.path+"/"+old_name)
					return
				}
			} else {
				log.Println(WARN, "unexpected error when checking if file exists before renaming:", err)
				return
			}
		}(dir, file)
	}
}

func DefaultTitle(seriesType string, namingScheme Option[string], path string, seasonPath string) string {
	var title string
	if seriesType == "singleSeasonNoMovies" || seriesType == "multipleSeasonNoMovies" || seriesType == "multipleSeasonWithMovies" {
		title = filepath.Base(path)
	} else if seriesType == "singleSeasonWithMovies" {
		title = filepath.Base(seasonPath)
	} else if seriesType == "namedSeasons" {
		title = filepath.Base(path) + " " + filepath.Base(seasonPath)
	}
	return CleanTitle(title)
}

func GenerateNewName(namingScheme Option[string], season_pad int, season_num int, ep_pad int, epNum int, title string, abs_path string) string {
	var newName string
	scheme, _ := namingScheme.Get()
	if namingScheme.IsSome() && scheme != "default" {
		// replace <season_num>
		newName = regexp.MustCompile(`<season_num(\s*:\s*\d+)?>`).ReplaceAllStringFunc(scheme, func(match string) string {
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
		newName = regexp.MustCompile(`<episode_num(\s*:\s*\d+)?>`).ReplaceAllStringFunc(newName, func(match string) string {
			// <episode_num: \d+>
			if strings.Contains(match, ":") {
				pad := regexp.MustCompile(`\d+`).FindString(match)
				pad_num, _ := strconv.Atoi(pad)
				return fmt.Sprintf("%0*d", pad_num, epNum)
			}
			// <episode_num>
			return fmt.Sprintf("%0*d", ep_pad, epNum)
		})
		// replace <self>
		newName = regexp.MustCompile(`<self\s*:\s*\d+,\d+>`).ReplaceAllStringFunc(newName, func(match string) string {
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
			return base_name[start : end+1]
		})
		// replace <parent> tokens with nth parent's name
		// lol goodluck: https://regex-vis.com/?r=%3C%28parent%28-parent%29*%28%5Cs*%3A%5Cs*%28%28%5Cd%2B%5Cs*%2C%5Cs*%5Cd%2B%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%7Cp%28-%5Cd%2B%29%3F%28%5Cs*%3A%5Cs*%28%28%5Cd%2B%5Cs*%2C%5Cs*%5Cd%2B%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%29%5Cs*%3E&e=0
		newName = regexp.MustCompile(`<(parent(-parent)*(\s*:\s*((\d+(\s*,\s*\d+)?)|('[^']*')))?|p(-\d+)?(\s*:\s*((\d+(\s*,\s*\d+)?)|('[^']*')))?)\s*>`).ReplaceAllStringFunc(newName, func(match string) string {
			n, err := ParentTokenToInt(match)
			if err != nil {
				return newName
			}
			parent_name := ParentN(abs_path, n)

			// <parent>
			if !strings.Contains(match, ":") {
				return parent_name
			}

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
				return parent_name[start : end+1]

			// <parent: '<regex_pattern>'>
			case '\'':
				regex_pattern := strings.Trim(val, "'")
				_, err := regexp.Compile(regex_pattern)
				if err != nil {
					return parent_name
				}
				sub_regexes := SplitRegexByPipe(regex_pattern)
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
				return parent_name[start : start+1]
			}
		})
		// append ext
		newName = filepath.Join(filepath.Dir(abs_path), fmt.Sprintf("%s%s", newName, filepath.Ext(abs_path)))

	} else if namingScheme.IsNone() || scheme == "default" {
		newName = fmt.Sprintf("S%0*dE%0*d %s%s",
			season_pad, season_num,
			ep_pad, epNum,
			title, filepath.Ext(abs_path))
		newName = filepath.Join(filepath.Dir(abs_path), newName)
	}

	return newName
}
