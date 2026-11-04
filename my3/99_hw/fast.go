package main

import (
	"bufio"
	"fmt"
	"hw3/easy"
	"io"
	"os"
	"slices"
	"strings"
	// "log"
)

const inputFile string = "./data/users.txt"

func FastSearch(out io.Writer) {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make([]string, 0, 100)
	uniqueBrowsers := 0
	foundUsers := strings.Builder{}

	reader := bufio.NewReader(file)
	var user easy.User
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		user.Browsers = user.Browsers[:0]
		// fmt.Printf("%v %v\n", err, line)
		err = user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}
		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {
			isContainsAndroid := strings.Contains(browserRaw, "Android")
			isContainsMSIE := strings.Contains(browserRaw, "MSIE")
			if isContainsAndroid || isContainsMSIE {
				idx := slices.Index(seenBrowsers, browserRaw)
				if isContainsAndroid {
					isAndroid = true
				}
				if isContainsMSIE {
					isMSIE = true
				}
				if idx == -1 {
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		// email := r.ReplaceAllString(string(user.Email), " [at] ")
		email := strings.Replace(user.Email, "@", " [at] ", 1)
		foundUsers.Write([]byte(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)))
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
