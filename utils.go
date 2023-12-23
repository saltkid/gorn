package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func is_media_file(file string) bool {
	// TODO: find a better way to identify media files
	media_extensions := map[string]bool {
		".mkv": true,
		".mp4": true,
		".avi": true,
		".mov": true,
		".webm": true,
		".ts": true,
	}
	return media_extensions[filepath.Ext(file)]
}

func has_movie (path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	seasonal_pattern := regexp.MustCompile(`^(?i)season\s+(\d+)`)
	specials_pattern := regexp.MustCompile(`^(?i)specials?|extras?|ova`)

	for _, file := range files {
		// found movie subdir
		if file.IsDir() && !seasonal_pattern.MatchString(file.Name()) && !specials_pattern.MatchString(file.Name()) {
			return true, nil
		}
	}

	// found no movie subdirs
	return false, nil
}


// filename renaming
type FilenameSort []string
// implement sort.Interface (Len, Less, Swap)

func (f FilenameSort) Len() int {
	return len(f)
}
func (f FilenameSort) Less(i, j int) bool {
	return compare_filenames(f[i], f[j])
}
func (f FilenameSort) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// compare two filenames which should come first based on numeric parts.
// if numeric parts are the same, compare non numeric parts character by character
func compare_filenames (f1 string, f2 string) bool {
	parts1 := split_filename(f1)
	parts2 := split_filename(f2)

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if is_numeric(parts1[i]) && is_numeric(parts2[i]) {
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
func split_filename (filename string) []string {
	var parts []string
	var current_part strings.Builder

	for i, c := range filename {
		// split when transitioning between numeric and non-numeric
		if i > 0 && (unicode.IsDigit(c) != unicode.IsDigit(rune(filename[i-1]))) {
			parts = append(parts, current_part.String())
			current_part.Reset()
		}
		// otherwise just add the character to the current part
		current_part.WriteRune(c)
	}
	parts = append(parts, current_part.String())

	return parts
}

func is_numeric (s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}