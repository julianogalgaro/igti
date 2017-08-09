package storage

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongo struct {
	ip   string
	port string
}

var session *mgo.Session

func (self *Mongo) getConnection() *mgo.Session {
	var mongoErr error

	for session == nil {
		session, mongoErr = mgo.Dial(self.ip)
		if mongoErr != nil {
			fmt.Println("Error connecting to mongo:", mongoErr)
			mongoErr = nil
			time.Sleep(2 * time.Second)
		}
	}
	session.Refresh()

	return session
}

func (self *Mongo) GetTweetWithoutClassification(t map[string]interface{}) error {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("repositorio").C("tweet")

	err := c.Find(bson.M{"classification": bson.M{"$exists": 0}}).One(&t)

	if err != nil {
		fmt.Println("Error" + err.Error())
		return err
	}
	return nil
}

func (self *Mongo) SetTweetClassification(id string, classification string) {
	s := self.getConnection()

	sessionCopy := s.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("repositorio").C("tweet")

	change := bson.M{"$set": bson.M{"classification": classification}}
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, change)
	if err != nil {
		fmt.Println("Activate Config Error:" + err.Error())
	}

	fmt.Println("TWEET [" + id + "] classificado como [" + classification + "]")
}

func NewMongo() *Mongo {

	m := &Mongo{
		ip:   "localhost",
		port: "27017",
	}

	return m
}
