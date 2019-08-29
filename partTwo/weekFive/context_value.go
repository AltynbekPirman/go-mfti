package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type key int

const timingsKey key = 3

func ctValue() {
	rand.Seed(time.Now().UTC().UnixNano())

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", loadPostsHandle)

	siteHandler := timingMiddleware(siteMux)

	err := http.ListenAndServe("127.0.0.1:9090", siteHandler)
	if err != nil {
		panic(err)
	}
}


func loadPostsHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	emulateWork(ctx, "check")
	emulateWork(ctx, "load")
	emulateWork(ctx, "load")
	time.Sleep(10 * time.Millisecond)
	emulateWork(ctx, "loadSidebar")
	emulateWork(ctx, "loadPics")
	_, err := fmt.Fprint(w, "done")
	if err != nil {
		log.Errorf("error, %s", err.Error())
	}
}

func timingMiddleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, timingsKey, &ctxTimings{
			Data: make(map[string]*Timing),
		})
		defer logContextTimings(ctx, r.URL.Path, time.Now())
		n.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logContextTimings(ctx context.Context, path string, t time.Time) {
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}
	total := time.Since(t)
	buf := bytes.NewBufferString(path)
	var tot time.Duration
	for timing, value := range timings.Data {
		tot += value.Duration
		buf.WriteString(fmt.Sprintf("\n\t%s(%d): %s", timing, value.Count, value.Duration))
	}
	buf.WriteString(fmt.Sprintf("\n\t total: %d", total))
	buf.WriteString(fmt.Sprintf("\n\t tracked: %d", tot))
	buf.WriteString(fmt.Sprintf("\n\t unkn: %d", total-tot))
	fmt.Println(buf.String())
}

func emulateWork(ctx context.Context, workName string) {
	defer trackContextTimeings(ctx, workName, time.Now())

	rnd := time.Duration(rand.Intn(50))
	time.Sleep(time.Millisecond * rnd)
}

type Timing struct {
	Count int
	Duration time.Duration
}

type ctxTimings struct {
	sync.Mutex
	Data map[string]*Timing
}

func trackContextTimeings(ctx context.Context, metricName string, t time.Time) {
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}

	elapsed := time.Since(t)
	timings.Lock()
	defer timings.Unlock()

	if metric, exists := timings.Data[metricName]; !exists {
		timings.Data[metricName] = &Timing{
			Count: 1,
			Duration: elapsed,
		}
	} else {
		metric.Count++
		metric.Duration += elapsed
	}
}