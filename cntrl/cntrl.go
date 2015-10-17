// cntrl package includes all web controllers
package cntrl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Maksadbek/go-url-shortener/models"
	"github.com/Maksadbek/go-url-shortener/mongostore"
	"github.com/gorilla/mux"
)

type ShortURLAPI struct {
	Mongo mongostore.MongoConn
}

var api ShortURLAPI

func init() {
	if err := api.Mongo.Connect(); err != nil {
		log.Println("can't connect to mongodb: ", err)
	}
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
		log.Println(shortURL)
	}
	longURL, err := api.Mongo.FindLongURL(shortURL)
	if err != nil {
		fmt.Fprintf(w, "could not find url")
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
