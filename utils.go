package main

import (
	"path/filepath"
	"os"
	"regexp"
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