package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)


var (
	h = testHandler{}
	ts = httptest.NewServer(h)
	scl = SearchClient{AccessToken: "test_token", URL: ts.URL}
)
const (
	NegativeRequestLimitText = "limit must be > 0"
	reqLimit = 25
	AccessToken = "test_token"
	BadToken = "Bad AccessToken"
)

type testHandler struct{

}

func (testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	token := headers["Accesstoken"]
	if len(token) == 0 {
		http.Error(w, "NoToken", http.StatusUnauthorized)
		return
	} else if token[0] != AccessToken{
		http.Error(w, "BadToken", http.StatusUnauthorized)
		return
	}

	users := []User{}
	users = append(users, User{1, "name", 25, "about", "gender"})
	jsonUsers, _ := json.Marshal(users)
	_, _ = fmt.Fprint(w, string(jsonUsers))
}


func TestAuth(t *testing.T) {
	cl := SearchClient{AccessToken: "tokenofbelisar", URL: ts.URL}
	sreq := SearchRequest{
		Limit: 20,

	}
	_, err := cl.FindUsers(sreq)
	if err != nil {
		if err.Error() != BadToken {
			t.Error("bad token")
		}
	}
}

func TestNegativeRequestLimit(t *testing.T){
	sreq := SearchRequest{
		Limit: -5,
	}
	sres, err := scl.FindUsers(sreq)
	if err != nil && sres == nil {
		if err.Error() != NegativeRequestLimitText {
			t.Errorf("Wrong Error text\n Expected: %s \n got: %s", NegativeRequestLimitText, err.Error())
		}
	} else {
		t.Error("Response is not nil for negative limit")
	}
}


func TestBigRequestLimit(t *testing.T) {
	sreq := SearchRequest{
		Limit: 5,
	}
	sres, err := scl.FindUsers(sreq)
	if sreq.Limit != reqLimit {
		fmt.Println(sres, err,  25)
		//t.Errorf("Request Limit not changed. Expected: %d, got: %d", reqLimit, sreq.Limit)
	}
}
