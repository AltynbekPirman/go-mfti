package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	return myDirTree(out, path, printFiles, 0, 0)
}

func myDirTree(out io.Writer, path string, printFiles bool, level int, dirLevel int) error {
	var pipeType string
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	fInfo, err := f.Readdir(0)
	if err != nil {
		return err
	}
	// remove files if not -f
	if !printFiles {
		fInfo = filterFiles(fInfo)
	}

	// sort files
	sort.Slice(fInfo, func(i, j int) bool {
		return fInfo[i].Name() < fInfo[j].Name()
	})

	fSize := len(fInfo)
	for ind, i := range fInfo {

		if ind == fSize - 1 {
			pipeType = "└───"
		} else {
			pipeType = "├───"
		}
		prePipes, err := drawPrePipes(level, dirLevel)
		if err != nil {
			return err
		}

		if i.IsDir() {
			txt := fmt.Sprintf("%s%s%s\n", prePipes, pipeType, i.Name())
			_, err := fmt.Fprint(out, txt)
			if err != nil {
				return err
			}
			if ind == fSize - 1 {
				// if file is in last directory do not change dirLevel to have empty tabs without pipes(|)
				err = myDirTree(out, strings.Join([]string{f.Name(), i.Name()}, "/"), printFiles, level + 1, dirLevel)
			} else {
				err = myDirTree(out, strings.Join([]string{f.Name(), i.Name()}, "/"), printFiles, level + 1, dirLevel+1)
			}
			if err != nil {
				return err
			}
		} else {
			if i.Size() != 0 {
				txt := fmt.Sprintf("%s%s%s (%db)\n", prePipes, pipeType, i.Name(), i.Size())
				_, err := fmt.Fprint(out, txt)
				if err != nil {
					return err
				}
			} else {
				txt := fmt.Sprintf("%s%s%s (%s)\n", prePipes, pipeType, i.Name(), "empty")
				_, err := fmt.Fprint(out, txt)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func filterFiles(fileSlice []os.FileInfo) []os.FileInfo {
	count := 0
	for _, i := range fileSlice {
		if i.IsDir() {
			fileSlice[count] = i
			count++
		}
	}
	return fileSlice[:count]
}

func drawPrePipes(n int, m int) (string, error) {

	if n == 0 {
		return "", nil
	}
	if m > n {
		return "", fmt.Errorf("invalid arguments n must be greater or equal to m")
	}

	var pipes []string
	for i := 0; i < m; i++ {
		pipes = append(pipes, "│\t")
	}

	for i := m; i < n; i++ {
		pipes = append(pipes, "\t")
	}

	res := strings.Join(pipes, "")
	return res, nil
}
