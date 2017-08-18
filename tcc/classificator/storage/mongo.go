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

func (self *Mongo) getClassificationDate(minute int) int64 {
	return time.Now().Add(time.Duration(minute) * time.Minute).UnixNano()
}

func (self *Mongo) GetTweetWithoutClassification(t map[string]interface{}) error {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()
	sessionCopy.Refresh()

	c := sessionCopy.DB("repositorio").C("tweet")

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"classificationDate": self.getClassificationDate(0)}},
		ReturnNew: false,
	}

	_, err := c.Find(bson.M{"classification": bson.M{"$exists": 0},
		"$or": []bson.M{
			bson.M{"classificationDate": bson.M{"$exists": 0}},
			bson.M{"classificationDate": bson.M{"$lt": self.getClassificationDate(-5)}},
		},
	}).Select(bson.M{"text": 1, "classificationPredict": 1, "_id": 1}).Limit(1).Apply(change, &t)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("Error: " + err.Error())
		return err
	}
	return nil
}

func (self *Mongo) GetTweetWithoutClassificationPredict(t map[string]interface{}) error {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()
	sessionCopy.Refresh()

	c := sessionCopy.DB("repositorio").C("tweet")

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"classificationPredictDate": self.getClassificationDate(0)}},
		ReturnNew: false,
	}

	_, err := c.Find(bson.M{"classificationPredict": bson.M{"$exists": 0},
		"$or": []bson.M{
			bson.M{"classificationPredictDate": bson.M{"$exists": 0}},
			bson.M{"classificationPredictDate": bson.M{"$lt": self.getClassificationDate(-1)}},
		},
	}).Select(bson.M{"text": 1, "_id": 1}).Limit(1).Apply(change, &t)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("Error: " + err.Error())
	}
	return err
}

func (self *Mongo) SetTweetClassification(id string, classification string) {
	s := self.getConnection()

	sessionCopy := s.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("repositorio").C("tweet")

	change := bson.M{"$set": bson.M{"classification": classification}}
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(id), "classification": bson.M{"$exists": 0}}, change)
	if err != nil {
		fmt.Println("ERROR on classify TWEET [" + id + "] as [" + classification + "]")
		return
	}

	fmt.Println("TWEET [" + id + "] classified success as [" + classification + "]")
}

func (self *Mongo) SetTweetClassificationPredict(id bson.ObjectId, classification string, rate float64) {
	s := self.getConnection()

	sessionCopy := s.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("repositorio").C("tweet")

	change := bson.M{"$set": bson.M{
		"classificationPredict":     classification,
		"classificationPredictRate": rate,
		"classificationPredictDate": self.getClassificationDate(0),
	}}
	err := c.Update(bson.M{"_id": id, "classificationPredict": bson.M{"$exists": 0}}, change)
	if err != nil {
		fmt.Println("ERROR on predict TWEET [" + id.Hex() + "] as [" + classification + "]")
		return
	}

	fmt.Println("TWEET [" + id.Hex() + "] predict as [" + classification + "]")
}

func (self *Mongo) GetCountClassification() (int, error) {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()
	sessionCopy.Refresh()

	c := sessionCopy.DB("repositorio").C("tweet")

	count, err := c.Find(bson.M{"classification": bson.M{"$exists": 1}}).Count()

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("Error: " + err.Error())
		return 0, err
	}
	return count, err
}

func (self *Mongo) GetCountClassificationPredict() (int, error) {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()
	sessionCopy.Refresh()

	c := sessionCopy.DB("repositorio").C("tweet")

	count, err := c.Find(bson.M{"classificationPredict": bson.M{"$exists": 1}}).Count()

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("Error: " + err.Error())
		return 0, err
	}
	return count, err
}

func (self *Mongo) GetCountTweets() (int, error) {
	s := self.getConnection()
	sessionCopy := s.Copy()
	defer sessionCopy.Close()
	sessionCopy.Refresh()

	c := sessionCopy.DB("repositorio").C("tweet")

	count, err := c.Find(nil).Count()

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("Error: " + err.Error())
		return 0, err
	}
	return count, err
}

func NewMongo() *Mongo {

	m := &Mongo{
		ip:   "localhost",
		port: "27017",
	}

	return m
}
