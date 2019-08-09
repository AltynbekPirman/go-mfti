package main

import (
	"fmt"
	"sync"
	"time"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	for _, j := range jobs {
		go j(in, out)
	}
	time.Sleep(time.Second*5)
}

func SingleHash(in, out chan interface{}) {
	for data := range in {
		fmt.Println("********************", data)
		a := DataSignerCrc32(fmt.Sprintf("%v", data))
		b := DataSignerCrc32(DataSignerMd5(fmt.Sprintf("%v", data)))
		out <- fmt.Sprintf(a, "~", b)
	}
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	var res string
	var multiHashes = make(map [int]string)
	for d := range in {
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				multiHashes[i] = DataSignerCrc32(fmt.Sprintf("%d%s", i, d))
				mu.Unlock()
			}()
		}
		wg.Wait()
		for i := 0; i < 6; i++ {
			res += multiHashes[i]
		}
		out <- res
	}
}

func CombineResults(in, out chan interface{}) {
	for i := range in {
		fmt.Println(i)
	}
	out<-"aaaaaaaaaaaa"
}
