// package mongostore can be used to work with MongoDB database 
package mongostore

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MONGODSN       string = "localhost"
	DBNAME         string = "urlshortener"
	COLLECTIONNAME string = "urlcollection"
)

var ErrMgoCreateColl = errors.New("could not create collection")
var ErrMgoConn = errors.New("could not connect to MongoDB server")
var ErrNilSession = errors.New("session is nil")
var ErrDupEntry = errors.New("duplicate entry")

// struct for keeping connection
type MongoConn struct {
	session *mgo.Session
}

// struct for mongo document, keeps each url description
type MongoDoc struct {
	ID       bson.ObjectId `bson:"_id"`
	ShortURL string        `bson:"shortURL"`
	LongURL  string        `bson:"longURL"`
}

// connect can be used to create a MongoDB connection
func (mc *MongoConn) Connect() (err error) {
	log.Println("connecting to MongoDB...")
	mc.session, err = mgo.Dial(MONGODSN)
	if err != nil {
		return ErrMgoConn
	}
	log.Println("successfully connected to MongoDB server")
	URLCollection := mc.session.DB(DBNAME).C(COLLECTIONNAME)
	if URLCollection == nil {
		err = ErrMgoCreateColl
	}
	index := mgo.Index{
		Key:      []string{"$text:shortURL"},
		Unique:   true,
		DropDups: true,
	}
	URLCollection.EnsureIndex(index)
	return n
}

// GetSessionAndCollection can be used as a connection pool, 
// on each call it copies and returns current session
func (mc *MongoConn) GetSessionAndCollection() (session *mgo.Session, collection *mgo.Collection, err error) {
	if mc.session != nil {
		session = mc.session.Copy()
		collection = session.DB(DBNAME).C(COLLECTIONNAME)
	} else {
		err = errors.New("session is nil")
	}
	return session, collection, err
}

// AddURLS can be used to insert url to MongoDB database
func (mc *MongoConn) AddURLs(longURL string, shortURL string) error {
	session, URLCollection, err := mc.GetSessionAndCollection()
	if err != nil {
		log.Println("could not get session and collection")
		return err
	}
	defer session.Close()
	err = URLCollection.Insert(
		&MongoDoc{
			ID:       bson.NewObjectId(),
			ShortURL: shortURL,
			LongURL:  longURL,
		},
	)
	if err != nil {
		// check for duplicate
		if mgo.IsDup(err) {
			return ErrDupEntry
		}
	}
	return nil
}

// FindLongURL can be used to find a url by its short alias
func (mc *MongoConn) FindLongURL(shortURL string) (string, error) {
	result := MongoDoc{}
	session, collection, err := mc.GetSessionAndCollection()
	if err != nil {
		return "", err
	}
	defer session.Close()
	err = collection.Find(bson.M{"shortURL": shortURL}).One(&result)
	if err != nil {
		return "", err
	}
	return result.LongURL, nil
}
