package main

import (
	"os"
	"path/filepath"
	"regexp"	
)

type Series struct {
	named_seasons               []string
	single_season_no_movies     []string
	single_season_with_movies   []string
	multiple_season_no_movies   []string
	multiple_season_with_movies []string
}

func (series *Series) split_series_by_type(series_entries []string) error {
	for _, series_entry := range series_entries {
		files, err := os.ReadDir(series_entry)
		if err != nil {
			return err
		}

		named_seasons_pattern := regexp.MustCompile(`^\d+\.\s+(.*)$`)
		seasonal_pattern := regexp.MustCompile(`^(?i)season\s+(\d+)`)

		for _, file := range files {
			if file.IsDir() {
				if file.Name() == filepath.Base(series_entry) {
					series.single_season_with_movies = append(series.single_season_with_movies, series_entry)
					break

				} else if named_seasons_pattern.MatchString(file.Name()) {
					series.named_seasons = append(series.named_seasons, series_entry)
					break

				} else if seasonal_pattern.MatchString(file.Name()) {
					has_movie, err := has_movie(series_entry)
					if err != nil {
						return err
					}

					if has_movie {
						series.multiple_season_with_movies = append(series.multiple_season_with_movies, series_entry)
						break

					} else {
						series.multiple_season_no_movies = append(series.multiple_season_no_movies, series_entry)
						break
					}
				}

			} else if is_media_file(file.Name()) {
				series.single_season_no_movies = append(series.single_season_no_movies, series_entry)
				break
			} 
		}
	}
	return nil
}