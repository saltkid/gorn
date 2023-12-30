package main

import (
	"os"
	"regexp"
)

type Movies struct {
	standalone []string
	movie_set  []string
}

func (movie *Movies) split_movies_by_type(movie_entries MovieEntries) error {
	for movie_entry := range movie_entries {
		files, err := os.ReadDir(movie_entry)
		if err != nil {
			return err
		}

		extras_pattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?`)

		for _, file := range files {
			if file.IsDir() && !extras_pattern.MatchString(file.Name()) {
				movie.movie_set = append(movie.movie_set, movie_entry)
				break

			} else if is_media_file(file.Name()) {
				movie.standalone = append(movie.standalone, movie_entry)
				break
			} 
		}
	}
	return nil
}