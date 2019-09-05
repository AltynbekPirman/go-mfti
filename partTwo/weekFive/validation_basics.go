package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"net/http"
)

type Message struct {
	ID 			int		`valid:",optional"`
	Priority 	string	`valid:"in(low|normal|high)"`
	Recipient 	string  `schema:"to" valid:"email"`
	Subject 	string	//`valid:"msgSubject"`
	Inner 		string  `schema:"-" valid:"-"`
	flag 		int
}


func handler1(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("request " + r.URL.String() + "\n\n"))
	//if err != nil {
	//	log.Error(err)
	//}
	msg := &Message{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	a := r.URL.Query()
	err := decoder.Decode(msg, a)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("Msg: %#v\n\n", msg)))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = govalidator.ValidateStruct(msg)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			for _, fld := range allErrs.Errors() {
				data := []byte(fmt.Sprintf("field: %#v\n\n", fld))
				_, err := w.Write(data)
				if err != nil {
					http.Error(w, "some err", http.StatusInternalServerError)
				}
			}
		}
		_, _ = w.Write([]byte(fmt.Sprintf("errors: %s\n\n", err)))
	} else {
		_, _ = w.Write([]byte("msg is correct"))
	}
}

func init() {
	govalidator.CustomTypeTagMap.Set("msgSubject",
		func(i interface{}, o interface{}) bool {
			subject, ok := i.(string)
			if !ok {
				return false
			}
			if len(subject) == 0 || len(subject) > 10 {
				return false
			}
			return true
		})
}

func validate() {
	http.HandleFunc("/", handler1)
	fmt.Println("starting server at :8080")
	_ = http.ListenAndServe(":9093", nil)
}
