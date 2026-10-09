package main

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

func getPrefix(path string, isLast bool) string {
	depth := strings.Count(path, string(os.PathSeparator))
	var prefix, suffix string
	if isLast {
		suffix = "└───"
	} else {
		suffix = "├───"
	}
	if isLast && depth != 0 {
		prefix = "│"
	}
	for range depth {
		prefix += "\t│"
	}
	prefix += suffix
	return prefix
}

func getSortedEntries(path string) ([]os.DirEntry, error) {
	d, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open dir %w", err)
	}

	entries, err := d.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read dir %w", err)
	}

	slices.SortFunc(entries, func(a os.DirEntry, b os.DirEntry) int {
		return cmp.Compare(a.Name(), b.Name())
	})

	return entries, nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	entries, err := getSortedEntries(path)
	if err != nil {
		return fmt.Errorf("getSortedEntries: %w", err)
	}

	for i, entry := range entries {
		prefix := getPrefix(path, i == len(entries)-1)
		entryName := entry.Name()
		if entry.IsDir() {
			fmt.Fprintf(out, "%s%s\n", prefix, entryName)
			dirTree(out, path+string(os.PathSeparator)+entryName, printFiles)
		} else if printFiles {
			entryInfo, err := entry.Info()
			if err != nil {
				return fmt.Errorf("Couldn't get file info %w", err)
			}

			fileSize := entryInfo.Size()
			if fileSize == 0 {
				fmt.Fprintf(out, "%s%s (%s)\n", prefix, entryName, "empty")
			} else {
				fmt.Fprintf(out, "%s%s (%db)\n", prefix, entryName, fileSize)
			}
		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
