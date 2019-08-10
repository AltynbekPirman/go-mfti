package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"runtime"
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
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
			close(out)
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
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
