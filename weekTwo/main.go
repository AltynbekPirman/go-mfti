package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sort"
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
	t := time.Now()
	ExecutePipeline(hashSignJobs...)
	fmt.Println("time: ", time.Since(t))
}


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


// 1173136728138862632818075107442090076184424490584241521304_1696913515191343735512658979631549563179965036907783101867
// _27225454331033649287118297354036464389062965355426795162684_29568666068035183841425683795340791879727309630931025356555
// _3994492081516972096677631278379039212655368881548151736_4958044192186797981418233587017209679042592862002427381542
// _4958044192186797981418233587017209679042592862002427381542