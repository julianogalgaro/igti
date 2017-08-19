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

func (self *Mongo) GetTweetsWithClassificationPredictByDate(limit int) ([]map[string]interface{}, error) {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()
	sessionCopy.Refresh()

	c := sessionCopy.DB("repositorio").C("tweet")

	startDate, err := self.GetLastProcessingDate("tweet")
	if err != nil {
		return nil, err
	}

	tweets := make([]map[string]interface{}, 0)

	err = c.Find(
		bson.M{"classificationPredict": bson.M{"$exists": 1},
			"classificationPredictDate": bson.M{"$gt": startDate},
		},
	).Sort("classificationPredictDate").Limit(limit).All(&tweets)
	fmt.Println("Total tweets found: ", len(tweets))
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("Error: " + err.Error())
		return nil, err
	}
	return tweets, nil
}

func (self *Mongo) SetLastProcessingDate(id string, lastDate int64) {
	s := self.getConnection()

	sessionCopy := s.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("controller").C("indexer")

	change := bson.M{"$set": bson.M{"lastProcessingDate": lastDate}}
	_, err := c.UpsertId(id, change)
	if err != nil {
		fmt.Println("ERROR on update lastProcessingDate to ", id)
		return
	}
}

func (self *Mongo) GetLastProcessingDate(id string) (int64, error) {
	s := self.getConnection()

	sessionCopy := s.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("controller").C("indexer")

	last := make(map[string]interface{})

	err := c.Find(
		bson.M{"_id": id},
	).One(&last)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("Error: " + err.Error())
		return int64(0), err
	}

	if err == mgo.ErrNotFound || last["lastProcessingDate"] == nil {
		self.SetLastProcessingDate(id, int64(0))
		fmt.Printf("[%s] actual lastProcessingDate [%d]\n", id, 0)
		return int64(0), nil
	} else {
		fmt.Printf("[%s] actual lastProcessingDate [%d]\n", id, last["lastProcessingDate"].(int64))
		return last["lastProcessingDate"].(int64), nil
	}
}

func NewMongo() *Mongo {

	m := &Mongo{
		ip:   "localhost",
		port: "27017",
	}

	return m
}
