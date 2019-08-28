package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func longSQLQuery() chan bool {
	ch := make(chan bool)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- true
	}()
	return ch
}

func multiproc() {
	timer := time.NewTimer(3*time.Second)
	select {
	case ct := <-timer.C:
		fmt.Println(ct)
	case <-time.After(time.Minute):
		fmt.Println("a minute passed")
	case res := <-longSQLQuery():

		// you can manually check and stop time.NewTimer() (and free resources), but you cannot stop time.After()
		var a time.Time
		if !timer.Stop() {
			a = <-timer.C
		}
		fmt.Println("got result", res, a)
	}

	// use ticker for periodic jobs
	ticker := time.NewTicker(time.Second)
	i := 0
	for tickTime := range ticker.C {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i > 5 {
			ticker.Stop() // stop does not close channel so you need to break from deadlock
			break
		}
	}
	fmt.Println("total", i)

	timerTwo := time.AfterFunc(5*time.Second, f)
	fmt.Scanln()
	timerTwo.Stop()
	fmt.Scanln()
}

func  f() {
	fmt.Println("after fuuuuuuuuuuuuuuuuunc")
}

// package context is used for cancellations (deadlines, timeouts) and to carry other request scoped values across APIs
func cntxt() {
	ctx, finish := context.WithCancel(context.Background()) // ctx's Done channel is closed after finish func is called
	result := make(chan int, 1)
	for i := 0; i <= 10; i++ {
		go worker(ctx, i, result)
	}
	foundBy := <-result
	fmt.Println("found by: ", foundBy)
	finish()
	time.Sleep(time.Second)

	workTime := 50*time.Millisecond
	ctx2, cancel := context.WithTimeout(context.Background(), workTime) //context is closed with timeout
	defer cancel() // releases resources if goroutine of context finishes faster
	res := make(chan int, 1)
	for i := 0; i <= 10; i++ {
		go worker(ctx2, i, res)
	}
	totalFound := 0
	LOOP:
		for {
			select {
			case <-ctx2.Done(): // returns context's channel
				break LOOP
			case foundBy := <- res:
				totalFound++
				fmt.Println("found by: ", foundBy)
			}
		}
	fmt.Println("total: ", totalFound)
	time.Sleep(time.Second)
}

func worker(ctx context.Context, n int, out chan<- int) {
	waitTime := time.Duration(rand.Intn(100)+10)*time.Millisecond
	fmt.Println(n, "sleep", waitTime)
	select {
	case <-ctx.Done():
		return
	case <-time.After(waitTime):
		fmt.Println("worker", n, "done")
		out <- n
	}
}

func getComments() chan string{
	comnts := make(chan string, 1)
	// async request to comments using goroutine and channels to pass data between goroutines
	go func(ch chan string) {
		fmt.Println("[getting comments]")
		time.Sleep(3*time.Second) // getting comments from DB)
		ch <- "all comments"
	}(comnts)
	return comnts
}

func getPage() {
	comntsCh := getComments()
	time.Sleep(3*time.Second) // getting article
	fmt.Println("My article")
	comnts := <- comntsCh
	fmt.Println(comnts)
}