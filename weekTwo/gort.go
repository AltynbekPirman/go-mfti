package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	iterationsNum = 7
	goroutinesNum = 6
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
	ch1 := make(chan int, 1) // chan with buffer 1 (chan can store one more value)
	go func(in chan int) {
		val := <- in
		fmt.Printf("received: %d\n", val)
		fmt.Println("*******")
	}(ch1)

	ch1 <- 42
	ch1 <- 42
	// Beware Deadlock! we try to write another value to chan with full buffer,
	// and all goroutines are dead (No one will read from this chan)
	//ch1 <- 42   --- fatal error: all goroutines are asleep - deadlock!

	in := make(chan int)

	// chan types: chan<- write only, <-chan read only, chan can be used both as read only and write only
	go func(out chan<- int) {
		for i := 0; i < 4; i++ {
			fmt.Println("before", i)
			out <- i
			fmt.Println("after", i)
		}
		//closing chan is useful when iterating over chan,
		// if chan is not closed and no one is filling it, iteration will cause DeadLock!
		close(out)
	}(in)
	for i := range in {
		fmt.Println("\treceived ", i)
	}
}

func multWSelect() {
	//	select can be used to multiplexing channels
	ch1 := make(chan string, 5)
	ch2 := make(chan string)
	go func() {
		val := <- ch2
		fmt.Println("received from ch2", val)
	}()
	iter := 0
	LOOP:
		for {
			select {
			case val := <-ch1:
				fmt.Println("ch1: ", val)
			case ch2 <- "aaaa":
				fmt.Println("pushed to ch2: ")
			default:
				time.Sleep(time.Second*2)
				fmt.Println(iter)
				iter++
				if iter > 2 {
					break LOOP // use label with break (otherwise break will work on select and loop will run forever)
				}
			}
		}
}

func gort() {
	for i := 0; i < goroutinesNum; i++ {
		go doSomeWork(i)
	}
	fmt.Scanln()
}
