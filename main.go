package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("skipping:", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		fmt.Println("walk failed:", err)
		os.Exit(1)
	}
}
