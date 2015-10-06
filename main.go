package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Maksadbek/shorturl/models"
	"github.com/Maksadbek/shorturl/mongostore"
	"github.com/gorilla/mux"
)

type ShortURLAPI struct {
	Mongo mongostore.MongoConn
}

var api ShortURLAPI

func main() {
	if err := api.Mongo.Connect(); err != nil {
		log.Println("Fuck, can't connect to mongo")
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", IndexHandler).Methods("GET")
	router.HandleFunc("/{shortURL}", RedirectHandler).Methods("GET")
	router.HandleFunc("/create", CreateHandler).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "simple url shortener service")
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	mapping := new(models.URLMapping)
	respEncoder := json.NewEncoder(w)

	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err = respEncoder.Encode(&models.APIResponse{StatusMsg: err.Error()}); err != nil {
			log.Println("could not encode response to json, ", err.Error())
		}
		return
	}
	if err := api.Mongo.AddURLs(mapping.LongURL, mapping.ShortURL); err != nil {
		w.WriteHeader(http.StatusConflict)
		if err = respEncoder.Encode(&models.APIResponse{StatusMsg: err.Error()}); err != nil {
			log.Println("could not encode response to json, ", err.Error())
		}
		return
	}
	respEncoder.Encode(&models.APIResponse{StatusMsg: "OK"})
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortURL, ok := params["shortURL"]
	if len(shortURL) == 0 && ok {
		fmt.Fprintf(w, "get the fuck off!")
		return
	}
	longURL, err := api.Mongo.FindLongURL(shortURL)
	if err != nil {
		fmt.Fprintf(w, "could not find url")
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
