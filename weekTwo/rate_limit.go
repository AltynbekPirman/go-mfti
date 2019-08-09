package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// we can limit number of active goroutines buffed chan, chan can consist of empty structs representing worker slots
func startWorkerRt(in int, wg *sync.WaitGroup, quotaCh chan struct{}) {
	quotaCh <- struct{}{} // write to chan, thus taking up slot or wait(block) if chan is full
	defer wg.Done()
	for j := 0; j < iterationsNum; j++ {
		time.Sleep(time.Second) // just to make results obvious
		fmt.Printf(formatWork(in, j))
	}
	<-quotaCh // read from chan, free slot for other goroutines
}


// this func is same as startWorkerRt, except it frees slot for other goroutines after each iteration
func rtJump(in int, wg *sync.WaitGroup, quotaCh chan struct{}) {
	quotaCh <- struct{}{}
	defer wg.Done()
	for j := 0; j < iterationsNum; j++ {
		time.Sleep(time.Second)
		fmt.Printf(formatWork(in, j))
		if j%2 == 0 {
			<-quotaCh
			quotaCh <- struct{}{}
		}
		runtime.Gosched()
	}
	<-quotaCh // read from chan, free slot for other goroutines
}

func rateLim() {
	wg := &sync.WaitGroup{}
	quotaCh := make(chan struct{}, 2) //buffer is 2, meaning only 2 goroutines can work simultaneously (try changing buff size and see output)
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		//go startWorkerRt(i, wg, quotaCh)
		go rtJump(i, wg, quotaCh)
	}
	wg.Wait()
}


// several goroutines writing and reading from one map leads to race condition.
// Note: running the program will not always result in error, making it difficult to find this type of bugs
// Use -race flag to check race condition ex: go run -race main.go
func raceCond(){
	counters := map[int]int{}
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{} // by locking and unlocking mutex before accessing data we can avoid race condition
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(counters map[int]int, th int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				mu.Lock()
				counters[th*10+j]++
				mu.Unlock()
			}
		}(counters, i)
	}
	wg.Wait()
	mu.Lock()
	fmt.Println(counters)
	mu.Unlock()
}

func cnAtomic() {
	var counter int32 = 0
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			//counter++ //will result in race condition and counter will be less than 1000
			atomic.AddInt32(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}