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
	LogArgs(args)

	err = start(args)
	if err != nil { log.Fatalln(FATAL, err) }
}

func start(args Args) error {
	defer timer("start")()
	var wg sync.WaitGroup
	defer wg.Wait() // for any early return errors

	seriesEntries, movieEntries, err := FetchEntries(args.root, args.series, args.movies)
	if err != nil { return err }

	// don't wait for logs to split entries by types
	go func(){
		LogRawEntries(seriesEntries, movieEntries)
		wg.Done()
	}()
	wg.Add(1)

	series := &Series{}
	err = series.SplitByType(seriesEntries)
	if err != nil { return err }
	
	// wait for previous log to finish printing to keep chronological order
	wg.Wait()
	go func(){
		series.LogEntries()
		wg.Done()
	}()
	wg.Add(1)

	// err = series.RenameEntries(args.options)
	// if err != nil { return err }

	movies := &Movies{}
	err = movies.SplitByType(movieEntries)
	if err != nil { return err }

	wg.Wait()
	go func() {
		movies.LogEntries()
		wg.Done()
	}()
	wg.Add(1)
	
	// err = movies.RenameEntries(args.options)
	// if err != nil { return err }

	return nil
}

func LogArgs(args Args) {
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
func FetchEntries(rootDirs []string, seriesDirs []string, movieDirs []string) ([]string, []string, error) {
	if len(rootDirs) == 0 && len(seriesDirs) == 0 && len(movieDirs) == 0 {
		return nil, nil, fmt.Errorf("passed no root, series, or movie directories")
	}

	entries := map[string][]string{
		"movies": make([]string, 0),
		"series": make([]string, 0),
	}
	for _, root := range rootDirs {
		separated, err := SeparateRoots(root)
		if err != nil {
			return nil, nil, err
		}

		for key, roots := range separated {
			for _, dir := range roots {
				subdirs, err := FetchSubdirs(dir)
				if err != nil {
					return nil, nil, err
				}
				entries[key] = append(entries[key], subdirs...)
			}
		}
	}

	for _, v := range seriesDirs {
		subdirs, err := FetchSubdirs(v)
		if err != nil {
			return nil, nil, err
		}
		entries["series"] = append(entries["series"], subdirs...)
	}
	for _, v := range movieDirs {
		subdirs, err := FetchSubdirs(v)
		if err != nil {
			return nil, nil, err
		}
		entries["movies"] = append(entries["movies"], subdirs...)
	}

	return entries["series"], entries["movies"], nil
}

func SeparateRoots(root string) (map[string][]string, error) {
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
		return nil, err
	}

	if len(rootDirs["movies"]) == 0 && len(rootDirs["series"]) == 0 {
		return nil, fmt.Errorf("no movie and series directory found")
	}

	return rootDirs, nil
}

func FetchSubdirs(dir string) ([]string, error) {
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
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries found under %s", dir)
	}

	return entries, nil
}
