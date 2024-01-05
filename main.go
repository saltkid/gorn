package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var version string

func main() {
	if len(os.Args) < 2 {
		WelcomeMsg(version)
		return
	}

	rawArgs, err := TokenizeArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}

	args, err := ParseArgs(rawArgs)
	if err != nil {
		if err.Error() != "safe exit" {
			panic(err)
		}
		return
	}

	if len(args.root) > 0 {
		fmt.Println("roots:")
		for _, root := range args.root {
			fmt.Println("\t", root)
		}
	}
	if len(args.series) > 0 {
		fmt.Println("series:")
		for _, series := range args.series {
			fmt.Println("\t", series)
		}
	}
	if len(args.movies) > 0 {
		fmt.Println("movies:")
		for _, movie := range args.movies {
			fmt.Println("\t", movie)
		}
	}
	ken, err := args.options.keepEpNums.Get()
	if err == nil {
		fmt.Println("keep episode numbers: ", ken)
	}
	sen, err := args.options.startingEpNum.Get()
	if err == nil {
		fmt.Println("starting episode number: ", sen)
	}
	ns, err := args.options.hasSeason0.Get()
	if err == nil {
		fmt.Println("naming scheme: ", ns)
	}

	seriesEntries, movieEntries, err := FetchEntries(args.root, args.series, args.movies)
	if err != nil {
		panic(err)
	}

	fmt.Println("series dirs (", len(seriesEntries), "): ")
	for _, series := range seriesEntries {
		fmt.Println("\t", series)
	}
	fmt.Println("movie dirs (", len(movieEntries), "): ")
	for _, movie := range movieEntries {
		fmt.Println("\t", movie)
	}
	fmt.Println()

	var series = Series{}
	err = series.SplitByType(seriesEntries)
	if err != nil {
		panic(err)
	}

	fmt.Println("categorized series: ")
	fmt.Println("namedSeasons: ")
	for _, v := range series.namedSeasons {
		fmt.Println("\t", v)
	}
	fmt.Println("singleSeasonNoMovies: ")
	for _, v := range series.singleSeasonNoMovies {
		fmt.Println("\t", v)
	}
	fmt.Println("singleSeasonWithMovies: ")
	for _, v := range series.singleSeasonWithMovies {
		fmt.Println("\t", v)
	}
	fmt.Println("multipleSeasonNoMovies: ")
	for _, v := range series.multipleSeasonNoMovies {
		fmt.Println("\t", v)
	}
	fmt.Println("multipleSeasonWithMovies: ")
	for _, v := range series.multipleSeasonWithMovies {
		fmt.Println("\t", v)
	}

	var movie = Movies{}
	err = movie.SplitByType(movieEntries)
	if err != nil {
		panic(err)
	}

	fmt.Println("categorized movies: ")
	fmt.Println("standalone: ")
	for _, v := range movie.standalone {
		fmt.Println("\t", v)
	}
	fmt.Println("movieSet: ")
	for _, v := range movie.movieSet {
		fmt.Println("\t", v)
	}

	fmt.Println("test for named seasons")
	namedSeasonOptions := PromptOptionalFlags(args.options, "all named seasons", 0)
	for _, v := range series.namedSeasons {
		info, err := SeriesRenamePrereqs(v, "namedSeasons", namedSeasonOptions)
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.Rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for single season no movies")
	ssnmOptions := PromptOptionalFlags(args.options, "all single season with no movies", 0)
	for _, v := range series.singleSeasonNoMovies {
		info, err := SeriesRenamePrereqs(v, "singleSeasonNoMovies", ssnmOptions)
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.Rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for single season with movies")
	sswmOptions := PromptOptionalFlags(args.options, "all single season with movies", 0)
	for _, v := range series.singleSeasonWithMovies {
		info, err := SeriesRenamePrereqs(v, "singleSeasonWithMovies", sswmOptions)
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.Rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for multiple season no movies")
	msnmOptions := PromptOptionalFlags(args.options, "all multiple season with no movies", 0)
	for _, v := range series.multipleSeasonNoMovies {
		info, err := SeriesRenamePrereqs(v, "multipleSeasonNoMovies", msnmOptions)
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.Rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for multiple season with movies")
	mswmOptions := PromptOptionalFlags(args.options, "all multiple season with movies", 0)
	for _, v := range series.multipleSeasonWithMovies {
		info, err := SeriesRenamePrereqs(v, "multipleSeasonWithMovies", mswmOptions)
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.Rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for standalone")
	for _, v := range movie.standalone {
		info, err := MovieRenamePrereqs(v, "standalone")
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.Rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for movie set")
	for _, v := range movie.movieSet {
		info, err := MovieRenamePrereqs(v, "movieSet")
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.Rename()
		if err != nil {
			panic(err)
		}
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
