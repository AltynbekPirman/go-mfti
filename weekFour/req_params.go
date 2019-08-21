package main

import (
	"fmt"
	"net/http"
	"reflect"
	"time"
)

func getParams() {

	paramMux := http.NewServeMux()

	server := http.Server{
		Addr: "127.0.0.1:8082",
		Handler: paramMux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	paramMux.HandleFunc("/", paramsHandler)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}


func paramsHandler(w http.ResponseWriter, r *http.Request) {
	// URL.Query() returns map[string][]string, it ignores malformed key value pairs
	p := r.URL.Query()
	fmt.Println(p)

	// FormValue returns value of passed key(if key is repeated - first occurrence is returned).
	// Returns empty string if no param is found
	k := r.FormValue("key")
	n := r.FormValue("name")
	// Do NOT do this!!! handle all errors. I did this to write this comment!!!
	_, _ = fmt.Fprint(w, reflect.TypeOf(k), k, n)

}


func login() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/logout", logoutPage)

	server := http.Server{
		Addr: "127.0.0.1:8083",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}


func mainPage(w http.ResponseWriter, r *http.Request) {
	session ,err := r.Cookie("session_id")
	var loggedIn bool
	if err != http.ErrNoCookie {
		loggedIn = true
	}
	w.Header().Set("RequestID", "1234567987897")
	_, err = fmt.Fprint(w, "<p><i><b>your browser:</b></i> " + r.UserAgent() + "<br>" + "<i><b>your accept:</b></i> " +
		r.Header.Get("Accept") +
		"</p><br><br>")
	if err != nil {
		fmt.Println(err)
	}

	if loggedIn {
		_, err = fmt.Fprint(w, `<a href="/logout">logout</a><br>`)
		_, err = fmt.Fprint(w, "<p>Welcome <b>", session.Value, "</b></p>")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = fmt.Fprint(w, `<a href="/login">login</a>`)
		if err != nil {
			fmt.Println(err)
		}
	}
}


func loginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(30 * time.Second)
	cookie := http.Cookie{
		Name: "session_id",
		Value: "belisar",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	session ,err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)

	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}