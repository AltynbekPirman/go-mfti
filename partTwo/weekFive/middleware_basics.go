package main

import (
	"fmt"
	"net/http"
	"time"
)

// we can manually handle all exceptions manually inside all handler functions.
// But it is better to use middleware instead of repeating the code
func pageWithAllChecks(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
	}()

	defer func(start time.Time) {
		fmt.Printf("[%s] %s %s %s\n", r.Method, r.RemoteAddr, r.URL, time.Since(start))
	}(time.Now())

	_, err := r.Cookie("session_id")

	if err != nil {
		fmt.Println("no auth", r.URL.Path)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// there will be your business logic
	panic("some_panic")
}


func handledFunc() {
	http.HandleFunc("/", pageWithAllChecks)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request){
		_, _ = w.Write([]byte("Login Page"))
	})

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}


// middleware example. Usually middleware can function that takes http Handler and returns another handler
func middleware() {
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/", adminIndex)

	// setting middleware
	adminHandler := adminAuthMiddleware(adminMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminHandler)

	siteMux.HandleFunc("/login", loginPage)

	siteMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		session, err := r.Cookie("session_id")
		// учебный пример! это не проверка авторизации!
		loggedIn := (err != http.ErrNoCookie)

		if loggedIn {
			fmt.Fprintln(w, "Welcome, "+session.Value)
		} else {
			fmt.Fprintln(w, `<a href="/login">login</a>`)
			fmt.Fprintln(w, "You need to login")
		}
	})

	h := accessLogMiddleware(siteMux)
	h = panicMiddleware(h)
	err := http.ListenAndServe("127.0.0.1:8000", h)
	if err != nil {
		panic(err)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(10 * time.Second)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   "belisar",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func adminIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Some important info in admin page")
}


func adminAuthMiddleware(next http.Handler) http.Handler {
	// use http.HandlerFunc
	// HandlerFunc(f) is a Handler that calls f(w, r).
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Println("auth Middleware", r.URL.Path)
		_, err := r.Cookie("session_id")

		if err != nil {
			fmt.Println("no auth", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func accessLogMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			fmt.Printf("[%s] %s %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(t))
		}(time.Now())
		next.ServeHTTP(w, r)
	})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panic middleware")
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered")
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

