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

type MongoConn struct {
	session *mgo.Session
}

type MongoDoc struct {
	ID       bson.ObjectId `bson:"_id"`
	ShortURL string        `bson:"shortURL"`
	LongURL  string        `bson:"longURL"`
}

func (mc *MongoConn) Connect() (err error) {
	log.Println("connecting to MongoDB...")
	mc.session, err = mgo.Dial(MONGODSN)
	if err != nil {
		return fmt.Errorf("could not connect to MongoDB server")
	}
	log.Println("successfully connected to MongoDB server")
	URLCollection := mc.session.DB(DBNAME).C(COLLECTIONNAME)
	if URLCollection == nil {
		err = errors.New("could not create collection")
	}
	index := mgo.Index{
		Key:      []string{"$text:shortURL"},
		Unique:   true,
		DropDups: true,
	}
	URLCollection.EnsureIndex(index)
	return nil
}

func (mc *MongoConn) GetSessionAndCollection() (session *mgo.Session, collection *mgo.Collection, err error) {
	if mc.session != nil {
		session = mc.session.Copy()
		collection = session.DB(DBNAME).C(COLLECTIONNAME)
	} else {
		err = errors.New("session is nil")
	}
	return session, collection, err
}

func (mc *MongoConn) AddURLs(longURL string, shortURL string) error {
	session, URLCollection, err := mc.GetSessionAndCollection()
	if err != nil {
		log.Println("could not get sessino and collection")
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
			return fmt.Errorf("duplicate entry for this short url %s", shortURL)
		}
	}
	return nil
}

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
