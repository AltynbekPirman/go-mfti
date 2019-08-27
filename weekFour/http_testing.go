package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func startServer() {
	mux := http.NewServeMux()

	server := http.Server{
		Addr: "127.0.0.1:8085",
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		_, _ = fmt.Fprintf(w, "request %#v\n", r)
		_, _ = fmt.Fprintf(w, "url %#v\n", r.URL)
	})
	mux.HandleFunc("/raw_body", func(w http.ResponseWriter, r *http.Request){
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		_, _ = fmt.Fprintf(w, "body: %s\n", string(body))
	})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}


func runGet() {
	url := "http://127.0.0.1:8085/?param=123&param2=test"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("get Body: %s\n", string(respBody))
}


func runGetFullReq(){
	req := &http.Request{
		Method: http.MethodGet,
		Header: http.Header{
			"User-Agent": {"application/golang"},
		},
	}
	var err error
	req.URL, err = url.Parse("http://127.0.0.1:8085/?id=45")
	if err != nil {
		fmt.Println(err)
		return
	}

	req.URL.Query().Set("user", "belisar")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("fullResp: %s\n", string(respBody))

}