package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


func main() {
	root, err := parse_args(os.Args)
	if err != nil {
		panic(err)
	}

	fmt.Println("root: ", root)

	root_dirs, err := get_root_dirs(root)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("root_dirs: ", root_dirs)

	entries, err := get_entries(root_dirs)
	if err != nil {
		panic(err)
	}

	fmt.Println("series_dirs: ", entries["series_dirs"])


	fmt.Println("movie_dirs: ", entries["movie_dirs"])
}

func parse_args(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("please add a root directory as an argument")
	}

	root, err := filepath.Abs(args[1])
	if err != nil {
		return "", err
	}
	return root, nil
}

func get_root_dirs(root string) (map[string]string, error) {
	root_dirs := map[string]string{
		"movie_dir":  "",
		"series_dir": "",
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
		// get only directories of depth 1 (directly under root)
		if d.IsDir() && path != root {
			parent_dir := filepath.Dir(path)
			if filepath.Base(parent_dir) == filepath.Base(root) {
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
			}
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if root_dirs["movie_dir"] == "" || root_dirs["series_dir"] == "" {
		return nil, fmt.Errorf("no movie or series directory found")
	}

	return root_dirs, nil
}

func get_entries(root_dirs map[string]string) (map[string][]string, error) {
	entries := map[string][]string{
		"movie_dirs": {},
		"series_dirs": {},
	}

	err := filepath.WalkDir(root_dirs["movie_dir"], func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// get only directories of depth 1 (directly under root)
		if d.IsDir() && path != root_dirs["movie_dir"] {
			entries["movie_dirs"] = append(entries["movie_dirs"], path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(root_dirs["series_dir"], func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// get only directories of depth 1 (directly under root)
		if d.IsDir() && path != root_dirs["series_dir"] {
			entries["series_dirs"] = append(entries["series_dirs"], path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}