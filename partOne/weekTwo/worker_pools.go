package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)


func workPool() {
	workerInput := make(chan string, 2)
	for i := 0; i < goroutinesNum; i++ {
		go startWorker(i, workerInput)
	}

	months := []string{"Jan", "Feb", "March", "April", "May",
				"June", "July", "Aug", "Sep", "Oct", "Nov", "Dec"}

	for _, monthName := range months {
		workerInput <- monthName
	}
	close(workerInput) //close chan so goroutines iterating over chan will know when chan ends
	time.Sleep(time.Millisecond) // keep main goroutine for some time (otherwise program will when close main goroutine ends)
}

func formatW(in int, j string) string  {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"recieved", j)
}

func startWorker(workerNum int, in <-chan string) {
	for input := range in {
		fmt.Printf(formatW(workerNum, input))
		runtime.Gosched()
	}
	printFinishWork(workerNum)
}

func printFinishWork(in int) {
	fmt.Println(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"===", in,
		"finished")
}

func workPoolWg() {
	wg := &sync.WaitGroup{}
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorkerWg(i, wg)
	}
	//time.Sleep(time.Millisecond)
	wg.Wait()	// blocks until wg counter is 0
}

func startWorkerWg(in int, wg *sync.WaitGroup) {
	defer wg.Done() // decreases wg counter by one, will panic if counter is negative
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		runtime.Gosched()
	}
}