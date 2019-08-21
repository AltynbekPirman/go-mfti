package main

import (
	"fmt"
	"net/http"
	"time"
)


func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Welcome home!\n")
	if err != nil {
		panic(err)
	}
	_, err = w.Write([]byte("Ahahahaha"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %s request from %s\n", r.Method, r.RemoteAddr)
}


// in case of using function handler, we cannot pass other arguments to function, thus
// there is no dependency injection or use of global variables. To solve this, you can use struct to serve content
// This struct must have ServeHTTP method
type Handler struct {
	Name string
}


func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "handler name: ", h.Name, " url: ", r.URL.String())
	if err != nil {
		panic(err)
	}
}


func httpServeWithFunc() {
	// handle pack takes in path and function that handles requests to the given path
	// handler func is strictly typed and must take http.ResponseWriter to write response and *http.Request as request struct
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request){
		_, err := fmt.Fprint(w, "Single page: ", r.URL.String())
		if err != nil {
			panic(err)
		}
	})
	http.HandleFunc("/pages/", func(w http.ResponseWriter, r *http.Request){
		_, err := fmt.Fprint(w, "Multiple pages: ", r.URL.String())
		if err != nil {
			panic(err)
		}
	})
	fmt.Println("starting server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}


func httpServeWithStruct() {
	h := &Handler{"New Handler"}
	h1 := &Handler{"home handler"}
	// takes pattern(url) string and struct with ServeHTTP method
	http.Handle("/", h)
	http.Handle("/home", h1)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}


func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, r.Body, r.Header)
	if err != nil {
		fmt.Println(err)
	}
}


func newHTTPServer() {

	// custom server can be created from new multiplexer. it gives more control over server configurations
	// multiple servers can be run by one program on multiple
	// goroutines (note: servers must run on different goroutines as listenAndServe is infinite loop)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	server := http.Server{
		Addr: "127.0.0.1:8081",
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}