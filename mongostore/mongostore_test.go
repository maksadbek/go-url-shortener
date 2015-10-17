package mongostore

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestConnection(t *testing.T) {
	var mc MongoConn
	err := mc.Connect()
	if err != nil {
		t.Error("can't connect to mongo: ", err)
	}
}

func TestGetSessionAndCollection(t *testing.T) {
	var mc MongoConn
	err := mc.Connect()
	if err != nil {
		t.Error(err)
	}
	_, _, err = mc.GetSessionAndCollection()
	if err != nil {
		t.Error(err)
	}
}

func TestAddURLs(t *testing.T) {
	var mc MongoConn
	var shortURL = "http://blablabla.corm"
	var longURL = "http://blasdkfalkdfalksdfja.rjwe"
	res := MongoDoc{}
	err := mc.Connect()
	if err != nil {
		t.Error(err)
	}
	// add new url
	err = mc.AddURLs(longURL, shortURL)
	if err != nil {
		t.Error(err)
	}
	// check if AddURLs worked fine
	session, collection, err := mc.GetSessionAndCollection()
	if err != nil {
		t.Error(err)
	}
	// close the session at the end
	defer session.Close()
	// find the url by shorturl
	err = collection.Find(bson.M{"shortURL": shortURL}).One(&res)
	if err != nil {
		t.Error(err)
	}
	// validate
	if got := res.LongURL; got != longURL {
		t.Errorf("want %s, got %s", longURL, got)
	}
	if got := res.ShortURL; got != shortURL {
		t.Errorf("want %s, got %s", shortURL, got)
	}
	// remove the test url
	err = collection.RemoveId(res.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestFindLongURL(t *testing.T) {
	var shortURL = "http://url"
	var longURL = "http://lurl"
	var mc MongoConn
	var res MongoDoc
	err := mc.Connect()
	if err != nil {
		t.Error(err)
	}
	// firstly, remove test data, to avoid duplicate entries
	session, collection, err := mc.GetSessionAndCollection()
	defer session.Close()
	err = collection.Find(bson.M{"shortURL": shortURL}).One(&res)
	// if shorturl found, then remove it by id
	if err != nil {
		t.Log("fuck")
	}
	if err == nil {
		err := collection.RemoveId(res.ID)
		if err != nil {
			t.Error(err)
		}
	}
	err = mc.AddURLs(longURL, shortURL)
	if err != nil {
		t.Error("error while adding new url", err)
	}
	got, err := mc.FindLongURL(shortURL)
	if err != nil {
		t.Error(err)
	}
	if got != longURL {
		t.Errorf("want %s, got %s", longURL, got)
	}
}
