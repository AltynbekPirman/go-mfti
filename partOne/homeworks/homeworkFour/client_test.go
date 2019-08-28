package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Note: tests are written without implementing SearchServer

const (
	NegativeRequestLimitText = "limit must be > 0"
	NegativeOffsetText = "offset must be > 0"
	reqLimit = 25
	AccessToken = "test_token"
	BadToken = "Bad AccessToken"
	BadServer = "SearchServer fatal error"
	ResponseTimeLimit = 2 * time.Second
)

var (
	h = testHandler{}
	ts = httptest.NewServer(h)
	scl = SearchClient{AccessToken: "test_token", URL: ts.URL}

)

type testHandler struct {

}

type tHandler struct {

}

type failedServerHandler struct {

}

type slowServerHandler struct {

}

type brokenServerHandler struct {

}

type invalidResponseHandler struct {

}

type badRequestHandler struct {

}

type badOrderFiledHandler struct {

}

type badErrorMessageHandler struct {

}

func (failedServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Internal Error", http.StatusInternalServerError)
}

func (badRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"Error": "some error"}`, http.StatusBadRequest)
}

func (badOrderFiledHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"Error": "ErrorBadOrderField"}`, http.StatusBadRequest)
}

func (badErrorMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "invalid error message", http.StatusBadRequest)
}

func (brokenServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}

func (slowServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(ResponseTimeLimit)
	_, _ = fmt.Fprint(w, "Some Response which shall not be available due to timeout")
}

func (invalidResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Not json response")
}

func (testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	params := r.URL.Query()
	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		http.Error(w, "limit is not number", 500)
		return
	}
	token := headers["Accesstoken"]
	if len(token) == 0 {
		http.Error(w, "NoToken", http.StatusUnauthorized)
		return
	} else if token[0] != AccessToken{
		http.Error(w, "BadToken", http.StatusUnauthorized)
		return
	}

	users := []User{
		{1, "name", 25, "about", "gender"}, {2, "name", 25, "about", "gender"},
		{3, "name", 25, "about", "gender"}, {4, "name", 25, "about", "gender"},
		{4, "name", 25, "about", "gender"}, {6, "name", 25, "about", "gender"},
		{7, "name", 25, "about", "gender"}, {8, "name", 25, "about", "gender"},
		{9, "name", 25, "about", "gender"}, {10, "name", 25, "about", "gender"},
		{11, "name", 25, "about", "gender"}, {12, "name", 25, "about", "gender"},
		{13, "name", 25, "about", "gender"}, {14, "name", 25, "about", "gender"},
		{15, "name", 25, "about", "gender"}, {16, "name", 25, "about", "gender"},
		{17, "name", 25, "about", "gender"}, {18, "name", 25, "about", "gender"},
		{19, "name", 25, "about", "gender"}, {20, "name", 25, "about", "gender"},
		{21, "name", 25, "about", "gender"}, {22, "name", 25, "about", "gender"},
		{23, "name", 25, "about", "gender"}, {24, "name", 25, "about", "gender"},
		{25, "name", 25, "about", "gender"}, {26, "name", 25, "about", "gender"},
		{27, "name", 25, "about", "gender"}, {28, "name", 25, "about", "gender"},
		{29, "name", 25, "about", "gender"}, {30, "name", 25, "about", "gender"},
		{31, "name", 25, "about", "gender"}, {32, "name", 25, "about", "gender"},
		{33, "name", 25, "about", "gender"}, {34, "name", 25, "about", "gender"},
		{35, "name", 25, "about", "gender"},
	}

	if limit < len(users) {
		jsonUsers, _ := json.Marshal(users[:limit])
		_, _ = fmt.Fprint(w, string(jsonUsers))

	} else {
		jsonUsers, _ := json.Marshal(users)
		_, _ = fmt.Fprint(w, string(jsonUsers))
	}
}


func (tHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	users := []User{
		{1, "name", 25, "about", "gender"}, {2, "name", 25, "about", "gender"},
		{3, "name", 25, "about", "gender"}, {4, "name", 25, "about", "gender"},
		{4, "name", 25, "about", "gender"}, {6, "name", 25, "about", "gender"},
		{7, "name", 25, "about", "gender"}, {8, "name", 25, "about", "gender"},
		{9, "name", 25, "about", "gender"}, {10, "name", 25, "about", "gender"},
	}
		jsonUsers, _ := json.Marshal(users)
		_, _ = fmt.Fprint(w, string(jsonUsers))
}


func TestAuth(t *testing.T) {
	http.NewServeMux()
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


func TestOffset(t *testing.T) {
	sreq := SearchRequest{
		Limit: 1,
		Offset: -5,
	}
	_, err := scl.FindUsers(sreq)
	if err.Error() != NegativeOffsetText {
		t.Errorf("offset error. Expected: %s, got: %s", NegativeOffsetText, err.Error())
	}
}

func TestBigRequestLimit(t *testing.T) {
	sreq := SearchRequest{
		Limit: 28,
		Offset: 0,
	}
	sres, _ := scl.FindUsers(sreq)
	if len(sres.Users) != reqLimit {
		t.Errorf("Request Limit not changed. Expected: %d, got: %d", reqLimit, len(sres.Users))
	}
}

func TestRequestLimit(t *testing.T) {
	h := tHandler{}
	ts := httptest.NewServer(h)
	defer ts.Close()
	sreq := SearchRequest{
		Limit: 10,
		Offset: 1,
	}
	cl := SearchClient{AccessToken: "test_token", URL: ts.URL}
	sres, _ := cl.FindUsers(sreq)
	if len(sres.Users) != sreq.Limit {
		t.Errorf("Request Limit was changed. Expected: %d, got: %d", reqLimit, sreq.Limit)
	}
}

func TestBadRequest(t *testing.T) {
	h := badRequestHandler{}
	ts := httptest.NewServer(h)
	defer ts.Close()
	sreq := SearchRequest{
		Limit: 10,
		Offset: 1,
	}
	cl := SearchClient{AccessToken: "test_token", URL: ts.URL}
	_, err := cl.FindUsers(sreq)
	if err == nil {
		t.Error("Timeout error expected, got nil for error")
	} else {
		if !strings.HasPrefix(err.Error(), "unknown bad request") {
			t.Error("Check error message")
		}
	}
}
func TestBadOrderField(t *testing.T) {
	h := badOrderFiledHandler{}
	ts := httptest.NewServer(h)
	defer ts.Close()
	sreq := SearchRequest{
		Limit: 10,
	}
	cl := SearchClient{AccessToken: "test_token", URL: ts.URL}
	_, err := cl.FindUsers(sreq)
	if err == nil {
		t.Error("Timeout error expected, got nil for error")
	} else {
		if !strings.HasPrefix(err.Error(), "Order") {
			t.Error("Check error message")
		}
	}
}
func TestBadErrorMessage(t *testing.T) {
	h := badErrorMessageHandler{}
	ts := httptest.NewServer(h)
	defer ts.Close()
	sreq := SearchRequest{
		Limit: 10,
	}
	cl := SearchClient{AccessToken: "test_token", URL: ts.URL}
	_, err := cl.FindUsers(sreq)
	if err == nil {
		t.Error("Timeout error expected, got nil for error")
	} else {
		if !strings.HasPrefix(err.Error(), "cant unpack") {
			t.Error("Check error message")
		}
	}
}

func TestSlowResponse(t *testing.T) {
	h := slowServerHandler{}
	ts := httptest.NewServer(h)
	defer ts.Close()
	sreq := SearchRequest{
		Limit: 10,
		Offset: 1,
	}
	cl := SearchClient{AccessToken: "test_token", URL: ts.URL}
	_, err := cl.FindUsers(sreq)

	if err == nil {
		t.Error("Timeout error expected, got nil for error")
	} else {
		if !strings.HasPrefix(err.Error(), "timeout") {
			t.Error("Check error message")
		}
	}
}

func TestBrokenServer(t *testing.T) {
	h := brokenServerHandler{}
	ts := httptest.NewServer(h)
	defer ts.Close()
	sreq := SearchRequest{
		Limit: 10,
		Offset: 1,
	}
	cl := SearchClient{AccessToken: "test_token", URL: ""}
	_, err := cl.FindUsers(sreq)

	if err == nil {
		t.Error("Timeout error expected, got nil for error")
	} else {
		if !strings.HasPrefix(err.Error(), "unknown error") {
			t.Error("Check error message")
		}
	}
}

func TestServerFail(t *testing.T) {
	fh := failedServerHandler{}
	fts := httptest.NewServer(fh)
	defer fts.Close()
	sreq := SearchRequest{
		Limit: 10,
		Offset: 1,
	}
	cl := SearchClient{AccessToken: "test_token", URL: fts.URL}
	_, err := cl.FindUsers(sreq)
	if err == nil {
		t.Error("expected error, got nil")
	} else if err.Error() != BadServer {
		t.Errorf("wrong error message. Expected: %s, got: %s", BadServer, err.Error() )
	}
}

func TestInvalidResponse(t *testing.T) {
	h := invalidResponseHandler{}
	ts := httptest.NewServer(h)
	defer ts.Close()
	sreq := SearchRequest{
		Limit: 10,
		Offset: 1,
	}
	cl := SearchClient{AccessToken: "test_token", URL: ts.URL}
	_, err := cl.FindUsers(sreq)
	if err == nil {
		t.Error("expected error, got nil")
	}
}