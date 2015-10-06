package mongostore

type MongoConn struct {
	session *mgo.Session
}

type MongoDoc struct {
	ID       bson.ObjectID `bson:"_id"`
	ShortURL string        `bson:"shortURL"`
	LongURL  string        `bson:"longURL"`
}

func (mc *MongoConnection) Connect() error {
	log.Println("connecting to MongoDB...")
	mc.session, err = mgo.Dial(MONGODSN)
	if err != nil {
		return fmt.Errorf("could not connect to MongoDB server")
	}
	log.Println("successfully connected to MongoDB server")
	URLCollection := c.session.DB(DBNAME).C(COLLECTIONNAME)
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

func (mc *MongoConnection) GetSessionAndCollection() (session *mgo.Session, collection *mgo.Collection, error err) {
	if mc.session != nil {
		session := mc.session.Copy()
		URLCollection = session.DB(DBNAME).C(COLLECTIONNAME)
	} else {
		err = errors.New("session is nil")
	}
	return
}

func (mc *MongoConnection) AddURLs(longURL string, shortURL string) error {
	session, URLCollection, err := mc.GetSessionAndCollection()
	if err != nil {
		log.Println("could not get sessino and collection")
		return err
	}
	defer session.Close()
	err = URLCollection.Insert(
		&MongoDOC{
			ID:       bson.NewObjectID(),
			ShortURL: shortURL,
			LongURL:  longURL,
		},
	)
	if err != nil {
		// check for duplicate
		if mgo.IsDup(err) {
			return errors.New("duplicate entry for this short url %s", shortURL)
		}
	}
	return nil
}
