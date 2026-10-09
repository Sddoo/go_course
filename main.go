package main

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"
)

func getSticks(isLast bool) string {
	if isLast {
		return "└───"
	} else {
		return "├───"
	}
}

func getSortedEntries(path string, printFiles bool) ([]os.DirEntry, error) {
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

	return slices.DeleteFunc(entries, func(entry os.DirEntry) bool {
		return !printFiles && !entry.IsDir()
	}), nil
}

func recursive(out io.Writer, path string, printFiles bool, prefix string) error {
	entries, err := getSortedEntries(path, printFiles)
	if err != nil {
		return fmt.Errorf("getSortedEntries: %w", err)
	}

	for i, entry := range entries {
		sticks := getSticks(i == len(entries)-1)
		entryName := entry.Name()
		if entry.IsDir() {
			var dirOffset string
			fmt.Fprintf(out, "%s%s\n", prefix+sticks, entryName)

			if i == len(entries)-1 {
				dirOffset = "\t"
			} else {
				dirOffset = "│\t"
			}

			if err := recursive(out, path+string(os.PathSeparator)+entryName, printFiles, prefix+dirOffset); err != nil {
				return err
			}
		} else if printFiles {
			entryInfo, err := entry.Info()
			if err != nil {
				return fmt.Errorf("Couldn't get file info %w", err)
			}

			fileSize := entryInfo.Size()
			if fileSize == 0 {
				fmt.Fprintf(out, "%s%s (%s)\n", prefix+sticks, entryName, "empty")
			} else {
				fmt.Fprintf(out, "%s%s (%db)\n", prefix+sticks, entryName, fileSize)
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return recursive(out, path, printFiles, "")
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
