// ===== PROGRESS NOTES =====
// DONE (Oct 28): Go set up on Pi, repo on GitHub,
//   baseline scanner: walks folders, reads lines, finds AWS keys,
//   skips node_modules and .git
//
// NEXT SESSION:
//   1. Check: `go run . testdata` prints ONLY [AWS KEY] testdata/config.py:1
//   2. Push test files: git add testdata && git commit -m "Add test fixtures" && git push
//      (if GitHub blocks it -> screenshot the message for review)
//   3. Exercise: add a Slack webhook regex, same pattern as AWS.
//      Should flag testdata/deploy/hook.yml
//
// AFTER THAT (concurrency, milestone 1):
//   Walker sends paths into a channel instead of calling scanFile directly.
//   Program must exit on its own (who closes the channel, and when?)
//
// REMEMBER: imports use "/", code uses "."; package things are Capitalized;
//   save (Ctrl+S) before go run; regex in `backticks`
// ==========================

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
