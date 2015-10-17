package main

import (
	"log"
	"net/http"

	"github.com/Maksadbek/go-url-shortener/cntrl"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", cntrl.IndexHandler).Methods("GET")
	router.HandleFunc("/{shortURL}", cntrl.RedirectHandler).Methods("GET")
	router.HandleFunc("/create", cntrl.CreateHandler).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
