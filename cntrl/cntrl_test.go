package cntrl

import (
	"bytes"
	"encoding/json"
	"github.com/Maksadbek/go-url-shortener/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// HOST name
const HOST string = "localhost:9999"

// TestCreateHandler tests short url creation
func TestCreateHandler(t *testing.T) {
	// marshal to json
	jsonMapping, err := json.Marshal(&mapping)
	if err != nil {
		t.Error(err)
	}
	// start new test server
	ts := httptest.NewServer(http.HandlerFunc(CreateHandler))
	defer ts.Close()
	// send a new post request
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(jsonMapping))
	// close the reader
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
	}
	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	var apiResp models.APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		t.Error(err)
	}
	if apiResp.StatusMsg != "OK" {
		t.Errorf("want %s, got %s", "OK", apiResp.StatusMsg)
	}
}

func TestRedirectHandler(t *testing.T) {
	// set up the router
	// otherwise, it can't parse the variable
	// see the explanation: https://groups.google.com/forum/#!topic/golang-nuts/Xs-Ho1feGyg
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{shortURL}", RedirectHandler).Methods("GET")
	// create new recorder
	w := httptest.NewRecorder()
	// create new test server
	req, err := http.NewRequest("GET", "http://localhost/a", nil)
	if err != nil {
		log.Fatal(err)
	}
	router.ServeHTTP(w, req)
	want := strings.TrimSpace(`<a href="` + mapping.LongURL + `">Found</a>.`)
	if got := strings.TrimSpace(w.Body.String()); got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}
