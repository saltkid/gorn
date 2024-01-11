package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var version string

func main() {
	defer timer("main")()

	if len(os.Args) < 2 {
		WelcomeMsg(version)
		return
	}
	rawArgs, err := TokenizeArgs(os.Args[1:])
	if err != nil { log.Fatalln(FATAL, err) }

	args, err := ParseArgs(rawArgs)
	if err != nil {
		// scuffed safe exit for --help and --version
		if _, ok := err.(SafeError); !ok { log.Fatalln(FATAL, err) }
		return
	}
	args.Log()

	seriesEntries, movieEntries := FetchEntries(args.root, args.series, args.movies)
	LogRawEntries(seriesEntries, movieEntries)
	
	var errChan = make(chan error, 2)
	wg := new(sync.WaitGroup)

	series := &Series{}
	wg.Add(1)
	go processMedia(series, seriesEntries, args.options, wg, errChan)

	movies := &Movies{}
	wg.Add(1)
	go processMedia(movies, movieEntries, args.options, wg, errChan)
	
	wg.Wait()
	close(errChan)
	for err := range errChan {
		log.Println(ERROR, err)
	}

	movies.LogEntries()
	series.LogEntries()
}

func processMedia(mediaFiles MediaFiles, entries []string, flags Flags, wg *sync.WaitGroup, errChan chan error) {
	defer timer("processMedia")()

	mediaFiles.SplitByType(entries)
	err := mediaFiles.RenameEntries(flags)
	if err != nil {
		errChan <- err
	}
	wg.Done()
}

func LogRawEntries(seriesEntries []string, movieEntries []string) {
	log.Println(INFO, "series dirs (", len(seriesEntries), "): ")
	for _, series := range seriesEntries {
		log.Println(INFO, "\t", series)
	}
	log.Println(INFO, "movie dirs (", len(movieEntries), "): ")
	for _, movie := range movieEntries {
		log.Println(INFO, "\t", movie)
	}
	fmt.Println()
}

// FetchEntries retrieves the series and movie entries from the given root, series, and movie directories.
//
// rootDirs: A slice of root directories to search for entries.
// seriesDirs: A slice of series directories to search for entries.
// movieDirs: A slice of movie directories to search for entries.
//
// Returns the series entries and movie entries as string slices.
func FetchEntries(rootDirs []string, seriesDirs []string, movieDirs []string) ([]string, []string) {
	entries := map[string][]string{
		"movies": make([]string, 0),
		"series": make([]string, 0),
	}
	for _, root := range rootDirs {
		separated := SeparateRoots(root)

		for key, roots := range separated {
			for _, dir := range roots {
				subdirs := FetchSubdirs(dir)
				entries[key] = append(entries[key], subdirs...)
			}
		}
	}

	for _, v := range seriesDirs {
		subdirs := FetchSubdirs(v)
		entries["series"] = append(entries["series"], subdirs...)
	}
	for _, v := range movieDirs {
		subdirs := FetchSubdirs(v)
		entries["movies"] = append(entries["movies"], subdirs...)
	}

	return entries["series"], entries["movies"]
}

func SeparateRoots(root string) (map[string][]string) {
	rootDirs := map[string][]string{
		"movies": {},
		"series": {},
	}
	validMoviePathNames := map[string]bool{
		"movies": true,
		"movie":  true,
	}
	validSeriesPathNames := map[string]bool{
		"series":  true,
		"shows":   true,
		"show":    true,
		"tv show": true,
		"tv":      true,
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Get only directories of depth 1 (directly under root)
			if path != root && filepath.Dir(path) == root {
				dirName := strings.ToLower(filepath.Base(path))
				if validMoviePathNames[dirName] {
					rootDirs["movies"] = append(rootDirs["movies"], path)
					return filepath.SkipDir

				} else if validSeriesPathNames[dirName] {
					rootDirs["series"] = append(rootDirs["series"], path)
					return filepath.SkipDir
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Println(ERROR, err)
		return map[string][]string{
			"movies": {},
			"series": {},
		}
	}

	if len(rootDirs["movies"]) == 0 && len(rootDirs["series"]) == 0 {
		log.Println(WARN, "no movie or series directory found under", root)
	}

	return rootDirs
}

func FetchSubdirs(dir string) []string {
	entries := []string{}
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Get only directories of depth 1 (directly under series dir) and does not start with a '.'
			if path != dir && filepath.Dir(path) == dir && !strings.HasPrefix(filepath.Base(path), ".") {
				entries = append(entries, path)
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		log.Println(ERROR, err)
		return []string{}
	}

	if len(entries) == 0 {
		log.Println(WARN, "no entries found under", dir)
	}

	return entries
}
