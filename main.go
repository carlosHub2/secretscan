package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// awsKeyPattern matches AWS access key IDs: "AKIA" + 16 uppercase lettersor digits.
var awsKeyPattern = regexp.MustCompile(`AKIA[A-Z0-9]{16}`)

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
			if d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		scanFile(path)
		return nil
	})
	if err != nil {
		fmt.Println("walk failed:", err)
		os.Exit(1)
	}
}

// scanFile opens one file, reads it line by line, and reports AWS keys.
func scanFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("cannot open:", path, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if awsKeyPattern.MatchString(line) {
			fmt.Printf("[AWS KEY] %s:%d\n", path, lineNum)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("read error:", path, err)
	}
}
