package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	// Plan
	// 1) Create scanner from file reader
	// 2) Scan file till line ends and unmarshal that json
	// 3) Get user and email from above json
	// 4) use sync pool for json. to optimize on allocation and gc of used json

	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	//seenBrowsers := []string{}
	seenBrowsers := make(map[string]bool)
	uniqueBrowsers := 0
	foundUsers := ""

	//var lines []string	// mem opt
	scanner := bufio.NewScanner(file)
	//users := make([]map[string]interface{}, 0)
	i := -1

	user := make(map[string]interface{})
	isAndroid := false
	isMime := false

	for scanner.Scan() {
		i++
		err := json.Unmarshal([]byte(scanner.Text()), &user)
		if err != nil {
			panic(err)
		}
		isAndroid = false
		isMime = false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				continue
			}
			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
				//notSeenBefore := true
				_, seen := seenBrowsers[browser]
				if !seen {
					seenBrowsers[browser] = true
					uniqueBrowsers++
				}
			}

			if ok := strings.Contains(browser, "MSIE"); ok {
				isMime = true
				_, seen := seenBrowsers[browser]
				if !seen {
					seenBrowsers[browser] = true
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMime) {
			// log.Println("cant cast browser to string")

			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)

	}

	//for i, user := range users {
	//
	//
	//}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}