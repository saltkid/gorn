// utils.go
// contains helper functions that are general purpose
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// checks if a file is a media file through checking the file extension
//
// current media extensions: .mkv, .mp4, .avi, .mov, .webm, .ts
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

// checks if a series entry contains a movie subdir through checking
// if the subdir name is not a season or an extras/specials subdir
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

// natural sorting of filenames
//
// usage: sort.Sort(FilenameSort(slice of paths))
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

// checks if a string is a valid number
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

// splits a regex by outermost pipes
//
// example:
//	`a|b|c|d|e` --> `a`, `b`, `c`, `d`, `e`
//	`a|(b|c|d|e)` --> `a`, `(b|c|d|e)`
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

// checks if a regex has only one match group
//
// note that this is usually called after SplitRegexByPipe so outermost pipes are removed first
//
// examples:
//	`(a|b)` --> true
//	`(a(b)a)` --> false
//	`(a|b|c)` --> true
//	`a` --> true
func HasOnlyOneMatchGroup(s string) bool {
	openingCount := 0
	closingCount := 0
	matchGroupCount := 0

	for _, c := range s {
		if c == '(' {
			openingCount++
		} else if c == ')' && closingCount < openingCount {
			closingCount++
			matchGroupCount++
		}
	}

	return matchGroupCount == 1 || matchGroupCount == 0
}

// returns the nth parent directory
//
// examples:
// 	ParentN("a/b/c", 1) --> "a/b"
// 	ParentN("a/b/c", 2) --> "a"
// 	ParentN("a/b/c", 3) --> ""
func ParentN(path string, n int) string {
	for i := 0; i < n; i++ {
		path = filepath.Dir(path)
	}
	return filepath.Base(path)
}

// specifically for exiting safely when user passed these switches:
//	`--help, -h`
//	`--version, -v`
type SafeError struct {
	safe error
}

// acts exactly like fmt.Errorf but returns a SafeError instead of error
func SafeErrorF(s string, args ...interface{}) SafeError {
	return SafeError{
		safe: fmt.Errorf(s, args...),
	}
}

// returns the underlying safe error as a string
func (s SafeError) Error() string {
	return s.safe.Error()
}

// Usage: put `defer timer("func_name")()` at the start of a function
//
// where "func_name" is just for logging purposes
//
// Reference:
//   - https://stackoverflow.com/questions/45766572/is-there-an-efficient-way-to-calculate-execution-time-in-golang
// thank you Cerise Limón
func timer(name string) func() {
    start := time.Now()
    return func() {
        log.Printf("%s %s took %v\n", TIME, name, time.Since(start))
    }
}

// for logging purposes
const (
	// for informational purposes
	INFO = "[INFO] " 					// no color

	// can safely skip error, doesn't interrupt process
	WARN = "\033[93m[WARN]\033[0m "		// yellow
	
	// cannot safely skip error, must interrupt process
	FATAL = "\033[91m[FATAL]\033[0m "	// red
	
	// for timing purposes
	TIME = "\033[94m[TIME]\033[0m "		// blue
)
