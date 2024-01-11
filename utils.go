package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func IsMediaFile(file string) bool {
	// TODO: find a better way to identify media files
	mediaExtensions := map[string]bool{
		".mkv":  true,
		".mp4":  true,
		".avi":  true,
		".mov":  true,
		".webm": true,
		".ts":   true,
	}
	return mediaExtensions[filepath.Ext(file)]
}

func HasMovie(path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	seasonalPattern := regexp.MustCompile(`^(?i)season\s+(\d+)`)
	extrasPattern := regexp.MustCompile(`(?i)^(specials?|extras?|o(n|v)a|trailers?|sub((title)?s)?|etc|others?)`)

	for _, file := range files {
		// found movie subdir
		if file.IsDir() && !seasonalPattern.MatchString(file.Name()) && !extrasPattern.MatchString(file.Name()) {
			return true, nil
		}
	}

	// found no movie subdirs
	return false, nil
}

// valid filename substring formats
//
// case insensitive
// can have spaces between season part and episode part (`S01 x E02`, `S01. E02`, `S01 _E02`, `S01 E02`)
// but can't have spaces between episode/season indicator and episode/season number.
// separators like `-` or `_` are allowed and can be repeated (`S01---E02`, `S01 __ E02`, `S01 xxE02`, `S01    E02`)
//
// S01E02 | S03.E04 | S05_E06 | S07-E08 | S09xE10 | S11 E12
//
// 01.02 | 03_04 | 05-06 | 07x08 | 09 10
//
// Episode 01 | Episode02 | EP03 | EP-04 | E_05 | EP.06
func ReadEpisodeNum(file string) (int, error) {

	// https://regex-vis.com/?r=s%5Cd%2B%5Cs*%28%3F%3Ax*%7C_*%7C-*%7C%5B.%5D*%29%5Cs*e%28%5Cd%2B%29%7C%5Cd%2B%5Cs*%28%3F%3Ax%2B%7C_%2B%7C-%2B%7C%5B.%5D%2B%29%5Cs*%28%5Cd%2B%29%7Cep%3F%28%3F%3Aisode%29%3F%5Cs*%28%3F%3A_%2B%7C-%2B%7C%5B.%5D%2B%29%5Cs*%28%5Cd%2B%29&e=0
	// match_id:											 			    [1]				   				   [2]									       [3]
	// captured:                                             				vv                    			   vv                 						   vv
	// substring:		                      s 01      x  _  -   .      e  02 | 03      x  _  -   .           04 |ep    isode          _  -   .           05
	episodePattern := regexp.MustCompile(`(?i)s\d+\s*(?:x*|_*|-*|[.]*)\s*e(\d+)|\d+\s*(?:x+|_+|-+|[.]+)\s*(\d+)|ep?(?:isode)?\s*(?:_+|-+|[.]+)\s*(\d+)`)
	match := episodePattern.FindStringSubmatch(file)
	if len(match) > 1 {
		epNumStr := ""
		for _, v := range match[1:] {
			if v != "" && epNumStr != "" {
				return 0, fmt.Errorf("multiple episode numbers found in %s: '%s', '%s' and '%s'", file, match[1], match[2], match[3])
			} else if v != "" {
				epNumStr = v
			}
		}
		if epNumStr == "" {
			return 0, fmt.Errorf("could not find episode number in %s", file)
		}

		epNum, err := strconv.Atoi(epNumStr)
		if err != nil {
			return 0, err
		}
		return epNum, nil
	} else {
		return 0, fmt.Errorf("could not find episode number in %s", file)
	}
}

// filename renaming
type FilenameSort []string

// implement sort.Interface (Len, Less, Swap)

func (f FilenameSort) Len() int {
	return len(f)
}
func (f FilenameSort) Less(i, j int) bool {
	return CompareFilenames(f[i], f[j])
}
func (f FilenameSort) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// compare two filenames which should come first based on numeric parts.
// if numeric parts are the same, compare non numeric parts character by character
func CompareFilenames(f1 string, f2 string) bool {
	parts1 := SplitFilename(f1)
	parts2 := SplitFilename(f2)

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if IsNumeric(parts1[i]) && IsNumeric(parts2[i]) {
			n1, _ := strconv.Atoi(parts1[i])
			n2, _ := strconv.Atoi(parts2[i])
			if n1 != n2 {
				return n1 < n2
			}
		} else if parts1[i] != parts2[i] {
			// loop through each char in non numeric parts and compare
			for j := 0; j < len(parts1[i]) && j < len(parts2[i]); j++ {
				if parts1[i][j] != parts2[i][j] {
					return parts1[i][j] < parts2[i][j]
				}
			}
		}
	}

	// all are the same, except one part is longer.
	return len(parts1) < len(parts2)
}

// split filename into numeric and non numeric parts for comparison purposes
func SplitFilename(filename string) []string {
	var parts []string
	var currentPart strings.Builder

	for i, c := range filename {
		// split when transitioning between numeric and non-numeric
		if i > 0 && (unicode.IsDigit(c) != unicode.IsDigit(rune(filename[i-1]))) {
			parts = append(parts, currentPart.String())
			currentPart.Reset()
		}
		// otherwise just add the character to the current part
		currentPart.WriteRune(c)
	}
	parts = append(parts, currentPart.String())

	return parts
}

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// remove the numbers in the title
//
//	"3. season 1 title" --> "season 1 title"
//	"04 - season 2 title" --> "season 2 title"
//	"005_season 3" --> "season 3"
//
// remove year in title; must be enclosed in parens
//
//	"title (2016)" --> "title"
//	"title (2000)" --> "title"
func CleanTitle(title string) string {
	// remove numbers
	re := regexp.MustCompile(`\d+\s*([.]|-|_)\s*`)
	title = re.ReplaceAllString(title, "")

	// remove year in parens
	re = regexp.MustCompile(`\s*\(\d{4}\)\s*`)
	title = re.ReplaceAllString(title, "")
	return title
}

func SplitRegexByPipe(s string) []string {
	var parts []string
	depth := 0
	part_start := 0

	for i, c := range s {
		if c == '|' && depth == 0 {
			parts = append(parts, s[part_start:i])
			part_start = i + 1
		} else if c == '(' {
			depth++
		} else if c == ')' {
			depth--
		}
	}

	parts = append(parts, s[part_start:])
	return parts
}

func HasOnlyOneMatchGroup(s string) bool {
	openingCount := 0
	closingCount := 0
	matchGroupCount := 0
	depth := 0

	for _, c := range s {
		if c == '(' && depth == 0 {
			openingCount++
			depth++
		} else if c == ')' && depth == 1 {
			closingCount++
			matchGroupCount++
			depth--
		}
	}

	return matchGroupCount == 1
}

func ParentTokenToInt(s string) (int, error) {
	// lmao https://regex-vis.com/?r=%3Cparent%28-parent%29*%28%5Cs*%3A%5Cs*%28%28%5Cd+%5Cs*%2C%5Cs*%5Cd+%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%5Cs*%3E&e=0
	longForm := regexp.MustCompile(`<parent(-parent)*(\s*:\s*((\d+(\s*,\s*\d+)?)|('[^']*')))?\s*>`)
	// lul https://regex-vis.com/?r=%3Cp%28-%5Cd%2B%29%3F%28%5Cs*%3A%5Cs*%28%28%5Cd%5Cs*%2C%5Cs*%5Cd%29%7C%28%27%5B%5E%27%5D*%27%29%29%29%3F%5Cs*%3E&e=0
	shortForm := regexp.MustCompile(`<p(-\d+)?(\s*:\s*((\d+(\s*,\s*\d+)?)|('[^']*')))?\s*>`)

	if longForm.MatchString(s) {
		return strings.Count(s, "parent"), nil

	} else if shortForm.MatchString(s) {
		// only p
		singleP := regexp.MustCompile(`p\s*[^-]:?`)
		if singleP.MatchString(s) {
			return 1, nil
		}
		// p-int
		matchNum := regexp.MustCompile(`p-(\d+)`).FindStringSubmatch(s)
		if len(matchNum) != 2 {
			return 0, fmt.Errorf("invalid parent token: %s", s)
		}
		num, err := strconv.Atoi(matchNum[1])
		if err != nil {
			return 0, err
		}
		return num, nil

	} else {
		return 0, fmt.Errorf("invalid parent token: %s", s)
	}
}

func ParentN(path string, n int) string {
	for i := 0; i < n; i++ {
		path = filepath.Dir(path)
	}
	return filepath.Base(path)
}
