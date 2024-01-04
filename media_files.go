package main

import (
	"os"
	"path/filepath"
	"regexp"
)

type MediaFiles interface {
	SplitByType(entries []string) error
}

type Movies struct {
	standalone []string
	movieSet   []string
}
type Series struct {
	namedSeasons             []string
	singleSeasonNoMovies     []string
	singleSeasonWithMovies   []string
	multipleSeasonNoMovies   []string
	multipleSeasonWithMovies []string
}

func (movie *Movies) SplitByType(movieEntries []string) error {
	for _, movieEntry := range movieEntries {
		files, err := os.ReadDir(movieEntry)
		if err != nil {
			return err
		}

		extrasPattern := regexp.MustCompile(`^(?i)specials?|extras?|trailers?`)

		for _, file := range files {
			if file.IsDir() && !extrasPattern.MatchString(file.Name()) {
				movie.movieSet = append(movie.movieSet, movieEntry)
				break

			} else if IsMediaFile(file.Name()) {
				movie.standalone = append(movie.standalone, movieEntry)
				break
			}
		}
	}
	return nil
}

func (series *Series) SplitByType(seriesEntries []string) error {
	for _, seriesEntry := range seriesEntries {
		files, err := os.ReadDir(seriesEntry)
		if err != nil {
			return err
		}

		namedSeasonsPattern := regexp.MustCompile(`^\d+\.\s+(.*)$`)
		seasonalPattern := regexp.MustCompile(`^(?i)season\s+(\d+)`)
		possiblySingleSeason := false
		for _, file := range files {
			if file.IsDir() {
				if file.Name() == filepath.Base(seriesEntry) {
					series.singleSeasonWithMovies = append(series.singleSeasonWithMovies, seriesEntry)
					possiblySingleSeason = false
					break

				} else if namedSeasonsPattern.MatchString(file.Name()) {
					series.namedSeasons = append(series.namedSeasons, seriesEntry)
					possiblySingleSeason = false
					break

				} else if seasonalPattern.MatchString(file.Name()) {
					HasMovie, err := HasMovie(seriesEntry)
					if err != nil {
						return err
					}

					if HasMovie {
						series.multipleSeasonWithMovies = append(series.multipleSeasonWithMovies, seriesEntry)
						possiblySingleSeason = false
						break

					} else {
						series.multipleSeasonNoMovies = append(series.multipleSeasonNoMovies, seriesEntry)
						possiblySingleSeason = false
						break
					}
				}

			} else if IsMediaFile(file.Name()) && !possiblySingleSeason {
				possiblySingleSeason = true
			}
		}

		if possiblySingleSeason {
			series.singleSeasonNoMovies = append(series.singleSeasonNoMovies, seriesEntry)
		}
	}
	return nil
}
