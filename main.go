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

	err = start(args)
	if err != nil { log.Fatalln(FATAL, err) }
}

func start(args Args) error {
	defer timer("start")()

	wg := new(sync.WaitGroup)
	defer wg.Wait() // for any early return errors

	go LogArgs(args, wg)

	seriesEntries, movieEntries := FetchEntries(args.root, args.series, args.movies)

	// don't wait for logs to split entries by types
	go LogRawEntries(seriesEntries, movieEntries, wg)

	series := &Series{}
	series.SplitByType(seriesEntries)
	
	// wait for previous log to finish printing to keep chronological order
	wg.Wait()
	go series.LogEntries(wg)

	// err = series.RenameEntries(args.options)
	// if err != nil { return err }

	movies := &Movies{}
	movies.SplitByType(movieEntries)

	wg.Wait()
	go movies.LogEntries(wg)
	
	// err = movies.RenameEntries(args.options)
	// if err != nil { return err }

	return nil
}

func LogArgs(args Args, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if len(args.root) > 0 {
		log.Println(INFO, "root directories: ")
		for _, root := range args.root {
			log.Println(INFO, "\t", root)
		}
	}
	if len(args.series) > 0 {
		log.Println(INFO, "series sources:")
		for _, series := range args.series {
			log.Println(INFO, "\t", series)
		}
	}
	if len(args.movies) > 0 {
		log.Println(INFO, "movies sources:")
		for _, movie := range args.movies {
			log.Println(INFO, "\t", movie)
		}
	}
	ken, err := args.options.keepEpNums.Get()
	if err == nil {
		log.Println(INFO, "keep episode numbers: ", ken)
	}
	sen, err := args.options.startingEpNum.Get()
	if err == nil {
		log.Println(INFO, "starting episode number: ", sen)
	}
	s0, err := args.options.hasSeason0.Get()
	if err == nil {
		log.Println(INFO, "has season 0: ", s0)
	}
	ns, err := args.options.namingScheme.Get()
	if err == nil {
		log.Println(INFO, "naming scheme: ", ns)
	}
}

func LogRawEntries(seriesEntries []string, movieEntries []string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

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
