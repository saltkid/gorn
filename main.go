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

	fmt.Println(root)

	root_dirs, err := get_root_dirs(root)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(root_dirs)
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