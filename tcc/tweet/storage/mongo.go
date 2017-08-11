package storage

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
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

func (self *Mongo) InsertTweet(t interface{}) {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()
	sessionCopy.Refresh()

	c := sessionCopy.DB("repositorio").C("tweet")

	err := c.Insert(t)

	if err != nil {
		panic(err)
	}
}

func NewMongo() *Mongo {

	m := &Mongo{
		ip:   "localhost",
		port: "27017",
	}

	return m
}
