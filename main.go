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
	for _, series := range series_entries {
		fmt.Println("\t", series)
	}
	fmt.Println("movie dirs (", len(movie_entries), "): ")
	for _, movie := range movie_entries {
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
	named_season_ken, named_season_sen, named_season_s0, named_season_ns := args.keep_ep_nums, args.starting_ep_num, args.has_season_0, args.naming_scheme
	prompt_additional_options(&named_season_ken, &named_season_sen, &named_season_s0, &named_season_ns, "named seasons")
	for _, v := range series.named_seasons {
		info, err := series_rename_prereqs(v, "named_seasons", named_season_ken, named_season_sen, named_season_s0, named_season_ns)
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
	// panic("test")

	fmt.Println("test for single season no movies")
	ssnm_ken, ssnm_sen, ssnm_s0, ssnm_ns := args.keep_ep_nums, args.starting_ep_num, args.has_season_0, args.naming_scheme
	prompt_additional_options(&ssnm_ken, &ssnm_sen, &ssnm_s0, &ssnm_ns, "single season no movies")
	for _, v := range series.single_season_no_movies {
		info, err := series_rename_prereqs(v, "single_season_no_movies", ssnm_ken, ssnm_sen, ssnm_s0, ssnm_ns)
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
	sswm_ken, sswm_sen, sswm_s0, sswm_ns := args.keep_ep_nums, args.starting_ep_num, args.has_season_0, args.naming_scheme
	prompt_additional_options(&sswm_ken, &sswm_sen, &sswm_s0, &sswm_ns, "single season with movies")
	for _, v := range series.single_season_with_movies {
		info, err := series_rename_prereqs(v, "single_season_with_movies", sswm_ken, sswm_sen, sswm_s0, sswm_ns)
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
	msnm_ken, msnm_sen, msnm_s0, msnm_ns := args.keep_ep_nums, args.starting_ep_num, args.has_season_0, args.naming_scheme
	prompt_additional_options(&msnm_ken, &msnm_sen, &msnm_s0, &msnm_ns, "multiple season no movies")
	for _, v := range series.multiple_season_no_movies {
		info, err := series_rename_prereqs(v, "multiple_season_no_movies", msnm_ken, msnm_sen, msnm_s0, msnm_ns)
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
	mswm_ken, mswm_sen, mswm_s0, mswm_ns := args.keep_ep_nums, args.starting_ep_num, args.has_season_0, args.naming_scheme
	prompt_additional_options(&mswm_ken, &mswm_sen, &mswm_s0, &mswm_ns, "multiple season with movies")
	for _, v := range series.multiple_season_with_movies {
		info, err := series_rename_prereqs(v, "multiple_season_with_movies", mswm_ken, mswm_sen, mswm_s0, mswm_ns)
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

// fetch_entries retrieves the series and movie entries from the given root, series, and movie directories.
//
// root_dirs: A slice of root directories to search for entries.
// series_dirs: A slice of series directories to search for entries.
// movie_dirs: A slice of movie directories to search for entries.
//
// Returns the series entries and movie entries as string slices.
func fetch_entries(root_dirs []string, series_dirs []string, movie_dirs []string) ([]string, []string, error) {
	if len(root_dirs) == 0 && len(series_dirs) == 0 && len(movie_dirs) == 0 {
		return nil, nil, fmt.Errorf("passed no root, series, or movie directories")
	}

	entries := map[string][]string{
		"movies":  make([]string, 0),
		"series": make([]string, 0),
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
				entries[key] = append(entries[key], subdirs...)
			}
		}
	}

	for _, v := range series_dirs {
		subdirs, err := fetch_subdirs(v)
		if err != nil {
			return nil, nil, err
		}
		entries["series"] = append(entries["series"], subdirs...)
	}
	for _, v := range movie_dirs {
		subdirs, err := fetch_subdirs(v)
		if err != nil {
			return nil, nil, err
		}
		entries["movies"] = append(entries["movies"], subdirs...)
	}

	return entries["series"], entries["movies"], nil
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