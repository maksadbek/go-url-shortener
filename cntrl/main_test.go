package cntrl

import (
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/Maksadbek/go-url-shortener/models"
)

// test url mapping data
var mapping = models.URLMapping{ShortURL: "abc", LongURL: "http://httpbin.org/get"}

func TestMain(m *testing.M) {
	retCode := m.Run()
	// remove test data from mongodb
	session, collection, err := api.Mongo.GetSessionAndCollection()
	defer session.Close()
	if err != nil {
		panic(err)
	}
	err = collection.Remove(bson.M{"shortURL": mapping.ShortURL})
	if err != nil {
		panic(err)
	}
	os.Exit(retCode)
}
