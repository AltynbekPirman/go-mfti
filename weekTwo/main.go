package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	//gort()
	//chans()
	//multWSelect()
	//multiproc()
	//cntxt()
	//getPage()
	//workPool()
	//workPoolWg()
	//rateLim()
	//raceCond()
	//cnAtomic()
	inputData := []int{1, 2}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				fmt.Println("job1: ", fibNum)
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			fmt.Println("end", data)
		}),
	}
	ExecutePipeline(hashSignJobs...)
}


func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})

	for _, j := range jobs {
		go j(in, out)
	}
	time.Sleep(time.Second*20)
}

func SingleHash(in, out chan interface{}) {
	var res, a, b string
	wg := &sync.WaitGroup{}
	for data := range out {
		fmt.Println(fmt.Sprintf("%v", data))
		wg.Add(2)
		go func(d interface{}){
			defer wg.Done()
			a = DataSignerCrc32(fmt.Sprintf("%s", data))
		}(data)
		go func(d interface{}){
			defer wg.Done()
			b = DataSignerMd5(fmt.Sprintf("%v", data))
		}(data)
		wg.Wait()
		b = DataSignerCrc32(b)
		res = fmt.Sprintf("single hash: %v~%v", a, b)
		fmt.Println(res)
		in <- res
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
			go func(ii int) {
				defer wg.Done()
				mu.Lock()
				multiHashes[ii] = DataSignerCrc32(fmt.Sprintf("%d%s", ii, d))
				mu.Unlock()
			}(i)
		}
		wg.Wait()
		for i := 0; i < 6; i++ {
			mu.Lock()
			res += multiHashes[i]
			mu.Unlock()
		}
		out <- res
	}
}


type job func(in, out chan interface{})


var (
	dataSignerOverheat uint32 = 0
	DataSignerSalt            = ""
)

var OverheatLock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
			fmt.Println("OverheatLock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var OverheatUnlock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
			fmt.Println("OverheatUnlock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var DataSignerMd5 = func(data string) string {
	OverheatLock()
	defer OverheatUnlock()
	data += DataSignerSalt
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond)
	return dataHash
}

var DataSignerCrc32 = func(data string) string {
	data += DataSignerSalt
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 10)
	time.Sleep(time.Second)
	return dataHash
}
