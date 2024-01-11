package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type MediaFiles interface {
	SplitByType(entries []string) error
	LogEntries()
	RenameEntries()
}

type Movies struct {
	standalone []string
	movieSet   []string
}
const (
	STANDALONE = "standalone"
	MOVIE_SET  = "movieSet"
)

type Series struct {
	namedSeasons             []string
	singleSeasonNoMovies     []string
	singleSeasonWithMovies   []string
	multipleSeasonNoMovies   []string
	multipleSeasonWithMovies []string
}
const (
	NAMED_SEASONS               = "namedSeasons"
	SINGLE_SEASON_NO_MOVIES     = "singleSeasonNoMovies"
	SINGLE_SEASON_WITH_MOVIES   = "singleSeasonWithMovies"
	MULTIPLE_SEASON_NO_MOVIES   = "multipleSeasonNoMovies"
	MULTIPLE_SEASON_WITH_MOVIES = "multipleSeasonWithMovies"
)

func (movie *Movies) SplitByType(entries []string) error {
	for _, movieEntry := range entries {
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

func (series *Series) SplitByType(entries []string) error {
	for _, seriesEntry := range entries {
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

func (movie *Movies) LogEntries() {
	log.Println(INFO, "categorized movies: ")
	log.Println(INFO, "standalone: ")
	for _, v := range movie.standalone {
		log.Println(INFO, "\t", v)
	}
	log.Println(INFO, "movie set: ")
	for _, v := range movie.movieSet {
		log.Println(INFO, "\t", v)
	}
}

func (series *Series) LogEntries() {
	log.Println(INFO, "categorized series: ")
	log.Println(INFO, "named seasons: ")
	for _, v := range series.namedSeasons {
		log.Println(INFO, "\t", v)
	}
	log.Println(INFO, "single season no movies: ")
	for _, v := range series.singleSeasonNoMovies {
		log.Println(INFO, "\t", v)
	}
	log.Println(INFO, "single season with movies: ")
	for _, v := range series.singleSeasonWithMovies {
		log.Println(INFO, "\t", v)
	}
	log.Println(INFO, "multiple season no movies: ")
	for _, v := range series.multipleSeasonNoMovies {
		log.Println(INFO, "\t", v)
	}
	log.Println(INFO, "multiple season with movies: ")
	for _, v := range series.multipleSeasonWithMovies {
		log.Println(INFO, "\t", v)
	}
}

func (movies *Movies) RenameEntries(options Flags) error {
	fmt.Println("Renaming standalone movies")
	for _, v := range movies.standalone {
		info, err := MovieRenamePrereqs(v, STANDALONE)
		if err != nil { return err }

		err = info.Rename()
		if err != nil { return err }
	}
	fmt.Println()

	fmt.Println("Renaming movie set")
	for _, v := range movies.movieSet {
		info, err := MovieRenamePrereqs(v, MOVIE_SET)
		if err != nil { return err }

		err = info.Rename()
		if err != nil { return err }
	}
	fmt.Println()

	return nil
}

func (series *Series) RenameEntries(options Flags) error {
	fmt.Println("Renaming named seasons")
	namedSeasonOptions := PromptOptionalFlags(options, "all named seasons", 0)
	for _, v := range series.namedSeasons {
		info, err := SeriesRenamePrereqs(v, NAMED_SEASONS, namedSeasonOptions)
		if err != nil {
			return err
		}

		err = info.Rename()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	fmt.Println("Renaming single season no movies")
	ssnmOptions := PromptOptionalFlags(options, "all single season with no movies", 0)
	for _, v := range series.singleSeasonNoMovies {
		info, err := SeriesRenamePrereqs(v, SINGLE_SEASON_NO_MOVIES, ssnmOptions)
		if err != nil {
			return err
		}

		err = info.Rename()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	fmt.Println("Renaming single season with movies")
	sswmOptions := PromptOptionalFlags(options, "all single season with movies", 0)
	for _, v := range series.singleSeasonWithMovies {
		info, err := SeriesRenamePrereqs(v, SINGLE_SEASON_WITH_MOVIES, sswmOptions)
		if err != nil {
			return err
		}

		err = info.Rename()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	fmt.Println("Renaming multiple season no movies")
	msnmOptions := PromptOptionalFlags(options, "all multiple season with no movies", 0)
	for _, v := range series.multipleSeasonNoMovies {
		info, err := SeriesRenamePrereqs(v, MULTIPLE_SEASON_NO_MOVIES, msnmOptions)
		if err != nil {
			return err
		}

		err = info.Rename()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	fmt.Println("Renaming multiple season with movies")
	mswmOptions := PromptOptionalFlags(options, "all multiple season with movies", 0)
	for _, v := range series.multipleSeasonWithMovies {
		info, err := SeriesRenamePrereqs(v, MULTIPLE_SEASON_WITH_MOVIES, mswmOptions)
		if err != nil {
			return err
		}

		err = info.Rename()
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}
