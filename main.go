package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var version string
var logLevel LogLevel = INFO_LEVEL // default log level

func main() {
	defer timer("main")()

	// handling input of user
	// errors can happen here and interrupt the process
	if len(os.Args) < 2 {
		WelcomeMsg(version)
		return
	}
	rawArgs, err := TokenizeArgs(os.Args[1:])
	if err != nil {
		gornLog(FATAL, err)
	}

	args, err := ParseArgs(rawArgs)
	if err != nil {
		// scuffed safe exit for --help and --version
		if _, ok := err.(SafeError); !ok {
			gornLog(FATAL, err)
		}
		return
	}
	args.Log()

	// renaming process
	// there should be no errors in renaming process since:
	//   - we already checked for errors in ParseArgs and TokenizeArgs
	//   - we can safely skip renaming a file if an error does occur
	seriesEntries, movieEntries := FetchEntries(args.root, args.series, args.movies)
	LogRawEntries(seriesEntries, movieEntries)

	wg := new(sync.WaitGroup)

	series := &Series{}
	wg.Add(1)
	go ProcessMedia(series, seriesEntries, args.options, wg)

	movies := &Movies{}
	wg.Add(1)
	go ProcessMedia(movies, movieEntries, args.options, wg)

	wg.Wait()
	movies.LogEntries()
	series.LogEntries()
}

// ProcessMedia first categorizes the entries by type, then renames them using the given flags.
//
//	mediaFiles: The media files to process. This should be a pointer to a struct that implements the MediaFiles interface
//	entries: The entries to categorize by type
//	flags: The flags to use for renaming
//	wg: The wait group to use for synchronization
func ProcessMedia(mediaFiles MediaFiles, entries []string, flags Flags, wg *sync.WaitGroup) {
	defer timer("ProcessMedia")()
	mediaFiles.SplitByType(entries)
	mediaFiles.RenameEntries(flags)
	wg.Done()
}

// LogRawEntries logs the uncategorized series and movie entries to the console.
func LogRawEntries(seriesEntries []string, movieEntries []string) {
	defer timer("LogRawEntries")()

	gornLog(INFO, "series dirs (", len(seriesEntries), "): ")
	for _, series := range seriesEntries {
		gornLog(INFO, "\t", series)
	}
	gornLog(INFO, "movie dirs (", len(movieEntries), "): ")
	for _, movie := range movieEntries {
		gornLog(INFO, "\t", movie)
	}
	fmt.Println()
}

// FetchEntries retrieves the series and movie entries from the given root, series, and movie directories.
//
//	rootDirs: A slice of root directories to search for entries.
//	seriesSourceDirs: A slice of series directories to search for entries.
//	movieSourceDirs: A slice of movie directories to search for entries.
//
//	Returns the series entries and movie entries as string slices.
func FetchEntries(rootDirs []string, seriesSourceDirs []string, movieSourceDirs []string) ([]string, []string) {
	defer timer("FetchEntries")()

	entries := map[string][]string{
		"movies": make([]string, 0),
		"series": make([]string, 0),
	}
	for _, root := range rootDirs {
		sources := FetchSourcesFromRoot(root)

		for key, roots := range sources {
			for _, dir := range roots {
				subdirs := FetchSubdirs(dir)
				entries[key] = append(entries[key], subdirs...)
			}
		}
	}

	for _, v := range seriesSourceDirs {
		subdirs := FetchSubdirs(v)
		entries["series"] = append(entries["series"], subdirs...)
	}
	for _, v := range movieSourceDirs {
		subdirs := FetchSubdirs(v)
		entries["movies"] = append(entries["movies"], subdirs...)
	}

	return entries["series"], entries["movies"]
}

// FetchSourcesFromRoot retrieves the source directories from the given root.
//
// Sources are directories containing entries:
//   - series source directories contains series entries
//   - movie source directories contains movie entries
func FetchSourcesFromRoot(root string) map[string][]string {
	sourceDirs := map[string][]string{
		"movies": {},
		"series": {},
	}
	validMovieSourceNames := map[string]bool{
		"movies": true,
		"movie":  true,
	}
	validSeriesSourceNames := map[string]bool{
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
				if validMovieSourceNames[dirName] {
					sourceDirs["movies"] = append(sourceDirs["movies"], path)
					return filepath.SkipDir

				} else if validSeriesSourceNames[dirName] {
					sourceDirs["series"] = append(sourceDirs["series"], path)
					return filepath.SkipDir
				}
			}
		}
		return nil
	})

	if err != nil {
		gornLog(WARN, "reading root directory error:", err)
	}

	if len(sourceDirs["movies"]) == 0 && len(sourceDirs["series"]) == 0 {
		gornLog(WARN, "no movie and/or series source directories found under:", root)
	}

	return sourceDirs
}

// FetchSubdirs retrieves the actual entries of the given source directory.
//   - if source is a series source, this returns series entries
//   - if source is a movie source, this returns movie entries
func FetchSubdirs(source string) []string {
	entries := []string{}
	err := filepath.WalkDir(source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Get only directories of depth 1 (directly under series dir) and does not start with a '.'
			if path != source && filepath.Dir(path) == source && !strings.HasPrefix(filepath.Base(path), ".") {
				entries = append(entries, path)
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		gornLog(WARN, "there was an error reading directory:", err)
	}

	if len(entries) == 0 {
		gornLog(WARN, "no entries found under", source)
	}

	return entries
}
