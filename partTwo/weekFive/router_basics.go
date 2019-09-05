package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func routers() {
	r := mux.NewRouter() // returns *Router, which implements http.Handler interface (http.ListenAndServe("9091", r))
	r.HandleFunc("/users", List).Host("localhost")
	r.HandleFunc("/users/{id:[0-9]+}/", getUser).Methods("get").Headers("Content-Type", "application/json")
	r.HandleFunc("/", List).Methods("GET", "Post") // list allowed methods

	err := http.ListenAndServe(":9091", r)
	if err != nil {
		panic(err)
	}
}

func mixRouters() {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/mux/user/{id:[0-9]+}", getUser).Methods("GET")
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", List)
	httpMux.Handle("/mux/", muxRouter)

	err := http.ListenAndServe(":9092", httpMux)
	panic(err)
}

func List(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "this is some list")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}


func getUser(w http.ResponseWriter, r *http.Request) {
	someUsers := map[string]string{"0": "Wrath", "1": "Pride", "2": "Greed",
		"3": "Lust",  "4": "Envy",
		"5": "Gluttony", "6": "Sloth",
	}
	vars := mux.Vars(r)
	id := vars["id"]

	if user, ok := someUsers[id]; ok {
		_, err := fmt.Fprint(w, "you meet: ", user)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "hmmmmm....", http.StatusBadRequest)
	}

}
