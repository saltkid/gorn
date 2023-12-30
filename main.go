package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args, err := parse_args(os.Args[1:])
	if err != nil {
		panic(err)
	}

	if len(args.root) > 0 {
		fmt.Println("root: ", args.root[0])
	}
	if len(args.series) > 0 {
		fmt.Println("series: ", args.series[0])
	}
	if len(args.movies) > 0 {
		fmt.Println("movies: ", args.movies[0])
	}
	ken, err := args.keep_ep_nums.get()
	if err == nil {
		fmt.Println("keep episode numbers: ", ken)
	}
	sen, err := args.starting_ep_num.get()
	if err == nil {
		fmt.Println("starting episode number: ", sen)
	}
	ns, err := args.naming_scheme.get()
	if err == nil {
		fmt.Println("naming scheme: ", ns)
	}

	series_entries, movie_entries, err := fetch_entries(args.root, args.series, args.movies)
	if err != nil {
		panic(err)
	}

	fmt.Println("series dirs (", len(series_entries), "): ")
	for series := range series_entries {
		fmt.Println("\t", series)
	}
	fmt.Println("movie dirs (", len(movie_entries), "): ")
	for movie := range movie_entries {
		fmt.Println("\t", movie)
	}
	fmt.Println()

	var series = Series{}
	err = series.split_by_type(series_entries)
	if err != nil {
		panic(err)
	}

	fmt.Println("categorized series: ")
	fmt.Println("named_seasons: ")
	for _, v := range series.named_seasons {
		fmt.Println("\t", v)
	}
	fmt.Println("single_season_no_movies: ")
	for _, v := range series.single_season_no_movies {
		fmt.Println("\t", v)
	}
	fmt.Println("single_season_with_movies: ")
	for _, v := range series.single_season_with_movies {
		fmt.Println("\t", v)
	}
	fmt.Println("multiple_season_no_movies: ")
	for _, v := range series.multiple_season_no_movies {
		fmt.Println("\t", v)
	}
	fmt.Println("multiple_season_with_movies: ")
	for _, v := range series.multiple_season_with_movies {
		fmt.Println("\t", v)
	}

	var movie = Movies{}
	err = movie.split_by_type(movie_entries)
	if err != nil {
		panic(err)
	}

	fmt.Println("categorized movies: ")
	fmt.Println("standalone: ")
	for _, v := range movie.standalone {
		fmt.Println("\t", v)
	}
	fmt.Println("movie_set: ")
	for _, v := range movie.movie_set {
		fmt.Println("\t", v)
	}

	fmt.Println("test for named seasons")
	for _, v := range series.named_seasons {
		info, err := series_rename_prereqs(v, "named_seasons", some[bool](false), some[int](1), some[bool](false))
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for single season no movies")
	for _, v := range series.single_season_no_movies {
		info, err := series_rename_prereqs(v, "single_season_no_movies", some[bool](false), some[int](1), some[bool](true))
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for single season with movies")
	for _, v := range series.single_season_with_movies {
		info, err := series_rename_prereqs(v, "single_season_with_movies", some[bool](true), some[int](1), some[bool](false))
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for multiple season no movies")
	for _, v := range series.multiple_season_no_movies {
		info, err := series_rename_prereqs(v, "multiple_season_no_movies", some[bool](false), some[int](1), some[bool](false))
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for multiple season with movies")
	for _, v := range series.multiple_season_with_movies {
		info, err := series_rename_prereqs(v, "multiple_season_with_movies", some[bool](false), some[int](1), some[bool](false))
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for standalone")
	for _, v := range movie.standalone {
		info, err := movie_rename_prereqs(v, "standalone")
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	fmt.Println("test for movie set")
	for _, v := range movie.movie_set {
		info, err := movie_rename_prereqs(v, "movie_set")
		if err != nil {
			panic(err)
		}
		fmt.Println(info)

		err = info.rename()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()
}

// set implementation for entries

type SeriesEntries map[string]struct{}
type MovieEntries map[string]struct{}
type Entries interface {
	Add(s ...string)
	Delete(s ...string)
	Has(s string) bool
}

func (se SeriesEntries) Add(s ...string) {
	for _, v := range s {
		se[v] = struct{}{}
	}
}
func (se SeriesEntries) Delete(s ...string) {
	for _, v := range s {
		delete(se, v)
	}
}
func (se SeriesEntries) Has(s string) bool {
	_, ok := se[s]
	return ok
}
func (me MovieEntries) Add(s ...string) {
	for _, v := range s {
		me[v] = struct{}{}
	}
}
func (me MovieEntries) Delete(s ...string) {
	for _, v := range s {
		delete(me, v)
	}
}
func (me MovieEntries) Has(s string) bool {
	_, ok := me[s]
	return ok
}

// fetch_entries retrieves the series and movie entries from the given root, series, and movie directories.
//
// root_dirs: A slice of root directories to search for entries.
// series_dirs: A slice of series directories to search for entries.
// movie_dirs: A slice of movie directories to search for entries.
//
// Returns the series entries and movie entries as SeriesEntries and MovieEntries respectively.
func fetch_entries(root_dirs []string, series_dirs []string, movie_dirs []string) (SeriesEntries, MovieEntries, error) {
	if len(root_dirs) == 0 && len(series_dirs) == 0 && len(movie_dirs) == 0 {
		return nil, nil, fmt.Errorf("passed no root, series, or movie directories")
	}

	entries := map[string]Entries{
		"movies":  MovieEntries{},
		"series": SeriesEntries{},
	}
	for _, root := range root_dirs {
		separated, err := separate_roots(root)
		if err != nil {
			return nil, nil, err
		}

		for key, roots := range separated {
			for _, dir := range roots {
				subdirs, err := fetch_subdirs(dir)
				if err != nil {
					return nil, nil, err
				}
				entries[key].Add(subdirs...)
			}
		}
	}

	for _, v := range series_dirs {
		subdirs, err := fetch_subdirs(v)
		if err != nil {
			return nil, nil, err
		}
		entries["series"].Add(subdirs...)
	}
	for _, v := range movie_dirs {
		subdirs, err := fetch_subdirs(v)
		if err != nil {
			return nil, nil, err
		}
		entries["movies"].Add(subdirs...)
	}

	return entries["series"].(SeriesEntries), entries["movies"].(MovieEntries), nil
}

func separate_roots(root string) (map[string][]string, error) {
	root_dirs := map[string][]string{
		"movies": {},
		"series": {},
	}
	valid_movie_path_names := map[string]bool{
		"movies": true,
		"movie":  true,
	}
	valid_series_path_names := map[string]bool{
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
			// get only directories of depth 1 (directly under root)
			if path != root && filepath.Dir(path) == root {
				dir_name := strings.ToLower(filepath.Base(path))
				if valid_movie_path_names[dir_name] {
					root_dirs["movies"] = append(root_dirs["movies"], path)
					return filepath.SkipDir

				} else if valid_series_path_names[dir_name] {
					root_dirs["series"] = append(root_dirs["series"], path)
					return filepath.SkipDir
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(root_dirs["movies"]) == 0 && len(root_dirs["series"]) == 0 {
		return nil, fmt.Errorf("no movie and series directory found")
	}

	return root_dirs, nil
}

func fetch_subdirs(dir string) ([]string, error) {
	entries := []string{}
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// get only directories of depth 1 (directly under series dir) and does not start with a '.'
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