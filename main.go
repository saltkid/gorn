package main

import (
	"fmt"
	"os"
	"path/filepath"
)


func main() {
	root, err := parse_args(os.Args)
	if err != nil {
		panic(err)
	}

	fmt.Println(root)
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