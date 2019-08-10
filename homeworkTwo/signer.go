package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	chans := make([]chan interface{}, 0)
	chans = append(chans, make(chan interface{}))
	for _, _ = range jobs{
		chans = append(chans, make(chan interface{}))
	}

	for i, j := range jobs {
		go j(chans[i], chans[i+1])
	}

	time.Sleep(time.Second*90)
}
func SingleHash(in, out chan interface{}) {
	defer func() {
		fmt.Println("closing singlehash")
		close(out)
	}()
	for data := range in {
		a := DataSignerCrc32(fmt.Sprintf("%v", data))
		b := DataSignerCrc32(DataSignerMd5(fmt.Sprintf("%v", data)))
		out <- fmt.Sprintf("%s~%s", a, b)
		runtime.Gosched()
		fmt.Println(fmt.Sprintf("%s~%s", a, b))
		runtime.Gosched()
	}
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	defer func() {
		fmt.Println("closing multihash")
		close(out)
	}()
	for d := range in {
		var res string
		var multiHashes = make(map [int]string)
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(ind int) {
				defer wg.Done()
				mu.Lock()
				multiHashes[ind] = DataSignerCrc32(fmt.Sprintf("%d%s", ind, d))
				mu.Unlock()
			}(i)
		}
		wg.Wait()
		for i := 0; i < 6; i++ {
			res += multiHashes[i]
		}
		out <- res
		fmt.Println("to combine: ", res)
		runtime.Gosched()
	}
}

func CombineResults(in, out chan interface{}) {
	var res string
	var r []string
	for i := range in {
		r = append(r, fmt.Sprintf("_%v", i))
	}
	for _, s := range r {
		res += s
	}
	out<-res[1:]

}