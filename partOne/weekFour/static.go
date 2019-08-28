package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)


func serveFile() {
	h := http.FileServer(http.Dir("./static"))

	mux := http.NewServeMux()
	mux.Handle("/", h)

	mux.HandleFunc("/form", formPage)
	mux.HandleFunc("/upload", uploadPage)
	mux.HandleFunc("/rr", uploadRawBody)

	server := http.Server{
		Addr: "127.0.0.1:8084",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}


var uploadFormTmpl = []byte(`
<html>
<body>
<form action="/upload" method="post" enctype="multipart/form-data">
Image: <input type="file" name="my_file">
<input type="submit" value="Upload">
</form>
</body>
</html>
`)


func formPage(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(uploadFormTmpl)
	if err != nil {
		panic(err)
	}
}


func uploadPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5*1024*1025)
	if err != nil {
		fmt.Println(err)
	}
	file, handl, err := r.FormFile("my_file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = fmt.Fprint(w, "handler.Filename %v\n", handl.Filename)
	_, err = fmt.Fprint(w, "handler.Filename %v\n", handl.Filename)
	if err != nil {
		fmt.Println(err)
	}

	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	_, err = fmt.Fprint(w, "md5 %x\n", hasher.Sum(nil))
	if err != nil {
		panic(err)
	}
}


type Params struct {
	ID string `json:"id"`
	Name string `json:"user"`
}

func uploadRawBody(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	defer r.Body.Close()

	p := &Params{}
	err = json.Unmarshal(body, p)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	_, err = fmt.Fprint(w, "content-type ", r.Header.Get("Content-Type"), "\n")
	_, err = fmt.Fprint(w, "params ", p, "\n")
	if err != nil {
		fmt.Println(err)
	}
}