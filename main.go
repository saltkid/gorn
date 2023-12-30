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
		fmt.Println("root: ", args.root[0].value)
	}
	if len(args.series) > 0 {
		fmt.Println("series: ", args.series[0].value)
	}
	if len(args.movies) > 0 {
		fmt.Println("movies: ", args.movies[0].value)
	}
	fmt.Println("keep episode numbers: ", args.keep_ep_nums.value)
	fmt.Println("starting episode number: ", args.starting_ep_num.value)
	fmt.Println("naming scheme: ", args.naming_scheme.value)

	// TODO: get entries differently for: root, series, movies
	var entries map[string][]string
	if len(args.root) > 0 {
		entries, err = get_root_dirs(args.root[0].value)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("series dirs (", len(entries["series_dirs"]), "): ")
	for _, series := range entries["series_dirs"] {
		fmt.Println("\t",series)
	}
	fmt.Println("movie dirs (", len(entries["movie_dirs"]), "): ")
	for _, movie := range entries["movie_dirs"] {
		fmt.Println("\t",movie)
	}
	fmt.Println()

	var series = Series{}
	err = series.split_series_by_type(entries["series_dirs"])
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
	err = movie.split_movies_by_type(entries["movie_dirs"])
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
		info, err := series_rename_prereqs(v, "named_seasons", false, 1, false)
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
		info, err := series_rename_prereqs(v, "single_season_no_movies", false, 1, true)
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
		info, err := series_rename_prereqs(v, "single_season_with_movies", true, 1, false)
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
		info, err := series_rename_prereqs(v, "multiple_season_no_movies", false, 1, false)
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
		info, err := series_rename_prereqs(v, "multiple_season_with_movies", false, 1, false)
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

func get_root_dirs(root string) (map[string][]string, error) {
	root_dirs := map[string]string{
		"movie_dir":  "",
		"series_dir": "",
	}
	entries := map[string][]string{
		"movie_dirs": {},
		"series_dirs": {},
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
			if path != root && filepath.Dir(path) == root && (root_dirs["movie_dir"] == "" || root_dirs["series_dir"] == "") {
				dir_name := strings.ToLower(filepath.Base(path))
				if valid_movie_path_names[dir_name] {
					if root_dirs["movie_dir"] == "" {
						root_dirs["movie_dir"] = path
					} else {
						return fmt.Errorf("multiple movie directories found")
					}

				} else if valid_series_path_names[dir_name] {
					if root_dirs["series_dir"] == "" {
						root_dirs["series_dir"] = path
					} else {
						return fmt.Errorf("multiple series directories found")
					}
				}

			// get only directories of depth 2 (directly under movie or series)
			} else if root_dirs["movie_dir"] != "" && filepath.Dir(path) == root_dirs["movie_dir"] {
				entries["movie_dirs"] = append(entries["movie_dirs"], path)
				return filepath.SkipDir
			} else if root_dirs["series_dir"] != "" && filepath.Dir(path) == root_dirs["series_dir"] {
				entries["series_dirs"] = append(entries["series_dirs"], path)
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if root_dirs["movie_dir"] == "" || root_dirs["series_dir"] == "" {
		return nil, fmt.Errorf("no movie or series directory found")
	}

	return entries, nil
}