package main

import (
	"os"
	"regexp"
	"path/filepath"
)

type MediaFiles interface {
	split_by_type(entries Entries) error
}

type Movies struct {
	standalone []string
	movie_set  []string
}
type Series struct {
	named_seasons               []string
	single_season_no_movies     []string
	single_season_with_movies   []string
	multiple_season_no_movies   []string
	multiple_season_with_movies []string
}

func (movie *Movies) split_by_type(movie_entries Entries) error {
	for movie_entry := range movie_entries.(MovieEntries) {
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

func (series *Series) split_by_type(series_entries Entries) error {
	for series_entry := range series_entries.(SeriesEntries) {
		files, err := os.ReadDir(series_entry)
		if err != nil {
			return err
		}

		named_seasons_pattern := regexp.MustCompile(`^\d+\.\s+(.*)$`)
		seasonal_pattern := regexp.MustCompile(`^(?i)season\s+(\d+)`)
		possibly_single_season := false
		for _, file := range files {
			if file.IsDir() {
				if file.Name() == filepath.Base(series_entry) {
					series.single_season_with_movies = append(series.single_season_with_movies, series_entry)
					possibly_single_season = false
					break

				} else if named_seasons_pattern.MatchString(file.Name()) {
					series.named_seasons = append(series.named_seasons, series_entry)
					possibly_single_season = false
					break

				} else if seasonal_pattern.MatchString(file.Name()) {
					has_movie, err := has_movie(series_entry)
					if err != nil {
						return err
					}

					if has_movie {
						series.multiple_season_with_movies = append(series.multiple_season_with_movies, series_entry)
						possibly_single_season = false
						break

					} else {
						series.multiple_season_no_movies = append(series.multiple_season_no_movies, series_entry)
						possibly_single_season = false
						break
					}
				}

			} else if is_media_file(file.Name()) && !possibly_single_season {
				possibly_single_season = true
			}
		}

		if possibly_single_season {
			series.single_season_no_movies = append(series.single_season_no_movies, series_entry)
		}
	}
	return nil
}