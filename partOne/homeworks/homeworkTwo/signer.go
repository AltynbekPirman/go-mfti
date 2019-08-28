package main

import (
	"fmt"
	"sort"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	chans := make([]chan interface{}, 0)
	chans = append(chans, make(chan interface{}))
	for _, _ = range jobs{
		chans = append(chans, make(chan interface{}))
	}
	wg := &sync.WaitGroup{}
	for i, j := range jobs {
		wg.Add(1)
		go func(a int, b job) {
			defer wg.Done()
			defer close(chans[a+1])
			b(chans[a], chans[a+1])
		}(i, j)
	}
	wg.Wait()
}

func hashCalc(mu *sync.Mutex, data interface{}, out chan interface{}) {
	var a, b string
	var md string
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		a = DataSignerCrc32(fmt.Sprintf("%v", data))
	}()
	go func() {
		defer wg.Done()
		mu.Lock()
		md = DataSignerMd5(fmt.Sprintf("%v", data))
		mu.Unlock()
		b = DataSignerCrc32(md)
	}()
	wg.Wait()
	out <- fmt.Sprintf("%s~%s", a, b)
}


func SingleHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	mu := &sync.Mutex{}
	for data := range in {
		wg.Add(1)
		go func(d interface{}, o chan interface{}){
			defer wg.Done()
			hashCalc(mu, d, o)
		}(data, out)
	}
	wg.Wait()
}


func multiHashCalc(w *sync.WaitGroup, out chan interface{}, d interface{}){
	defer w.Done()
	wg := &sync.WaitGroup{}
	resChan := make(chan string, 6)
	defer close(resChan)
	var res string
	r := make([]string, 0)
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(ind int) {
			defer wg.Done()
			m := DataSignerCrc32(fmt.Sprintf("%d%s", ind, d))
			resChan <- fmt.Sprintf("%d%s", ind, m)
		}(i)
	}
	wg.Wait()
	for i := 0; i < 6; i++ {
		r = append(r, <-resChan)
	}
	sort.Strings(r)
	for _, i := range r {
		res += i[1:]
	}
	out <- res
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for d := range in {
		wg.Add(1)
		go multiHashCalc(wg, out, d)
	}
	wg.Wait()

}

func CombineResults(in, out chan interface{}) {
	var res string
	var r []string
	for i := range in {
		r = append(r, fmt.Sprintf("_%v", i))
	}
	sort.Strings(r)
	for _, s := range r {
		res += s
	}
	out<-res[1:]
}