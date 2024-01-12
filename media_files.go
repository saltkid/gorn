package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type MediaFiles interface {
	SplitByType(entries []string)
	RenameEntries(Flags)
	LogEntries()
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

func (movie *Movies) SplitByType(entries []string) {
	for _, movieEntry := range entries {
		files, err := os.ReadDir(movieEntry)
		if err != nil {
			log.Println(WARN, "there was anerror reading movie entry:", err, "; skipping categorizing entry:", movieEntry)
			continue
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
}

func (series *Series) SplitByType(entries []string) {
	for _, seriesEntry := range entries {
		files, err := os.ReadDir(seriesEntry)
		if err != nil {
			log.Println(WARN, "there was an error reading series entry:", err, "; skipping categorizing entry:", seriesEntry)
			continue
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
					HasMovie, _ := HasMovie(seriesEntry)

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
}

func (movies *Movies) RenameEntries(options Flags) {
	wg := new(sync.WaitGroup)

	log.Println(INFO, "Renaming standalone movies..")
	for _, v := range movies.standalone {
		wg.Add(1)
		go func(v string){
			info := MovieRenamePrereqs(v, STANDALONE)
			info.Rename(wg)
			wg.Done()
		}(v)
	}
	log.Println(INFO, "Renaming movie sets...")
	for _, v := range movies.movieSet {
		wg.Add(1)
		go func(v string){
			info := MovieRenamePrereqs(v, MOVIE_SET)
			info.Rename(wg)
			wg.Done()
		}(v)
	}

	wg.Wait()
	log.Println(INFO, "Done renaming movie entries.")
}

func (series *Series) RenameEntries(options Flags) {
	wg := new(sync.WaitGroup)

	log.Println(INFO, "Renaming named seasons...")
	namedSeasonOptions := PromptOptionalFlags(options, "all named seasons", 0)
	for _, v := range series.namedSeasons {
		wg.Add(1)
		go func(v string){
			info := SeriesRenamePrereqs(v, NAMED_SEASONS, namedSeasonOptions)
			info.Rename()
			wg.Done()
		}(v)
	}
	
	log.Println(INFO, "Renaming single season no movies...")
	ssnmOptions := PromptOptionalFlags(options, "all single season with no movies", 0)
	for _, v := range series.singleSeasonNoMovies {
		wg.Add(1)
		go func(v string){
			info := SeriesRenamePrereqs(v, SINGLE_SEASON_NO_MOVIES, ssnmOptions)
			info.Rename()
			wg.Done()
		}(v)
	}
	
	log.Println(INFO, "Renaming single season with movies...")
	sswmOptions := PromptOptionalFlags(options, "all single season with movies", 0)
	for _, v := range series.singleSeasonWithMovies {
		wg.Add(1)
		go func(v string) {
			info := SeriesRenamePrereqs(v, SINGLE_SEASON_WITH_MOVIES, sswmOptions)
			info.Rename()
			wg.Done()
		}(v)
	}
	
	log.Println(INFO, "Renaming multiple season no movies...")
	msnmOptions := PromptOptionalFlags(options, "all multiple season with no movies", 0)
	for _, v := range series.multipleSeasonNoMovies {
		wg.Add(1)
		go func(v string){
			info := SeriesRenamePrereqs(v, MULTIPLE_SEASON_NO_MOVIES, msnmOptions)
			info.Rename()
			wg.Done()
		}(v)
	}
	
	log.Println(INFO, "Renaming multiple season with movies...")
	mswmOptions := PromptOptionalFlags(options, "all multiple season with movies", 0)
	for _, v := range series.multipleSeasonWithMovies {
		wg.Add(1)
		go func(v string){
			info := SeriesRenamePrereqs(v, MULTIPLE_SEASON_WITH_MOVIES, mswmOptions)
			info.Rename()
			wg.Done()
		}(v)
	}
	
	wg.Wait()
	log.Println(INFO, "Done renaming series entries.")
}

func (movie *Movies) LogEntries() {
	defer timer("movies LogEntries")()

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
	defer timer("series LogEntries")()
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
