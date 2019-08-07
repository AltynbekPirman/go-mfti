package main

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func doSomeWork(in int) {
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		runtime.Gosched() // frees processor to run other goroutines, does not suspend goroutine -> it will auto-resume
	}
}

func formatWork(in, j int) string  {
	return fmt.Sprintln(
		strings.Repeat(" ", in), "*",
		strings.Repeat(" ", goroutinesNum - in), "th", in,
		"iter", j, strings.Repeat("*", j))
}

// channels are used to pass data between channels
func chans() {

}

func gort() {
	for i := 0; i < goroutinesNum; i++ {
		go doSomeWork(i)
	}
	fmt.Scanln()
}
