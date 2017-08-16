package storage

import "gopkg.in/mgo.v2/bson"

type Storage interface {
	GetTweetWithoutClassification(map[string]interface{}) error
	SetTweetClassification(string, string)
	GetCountClassification() (int, error)
	GetCountTweets() (int, error)
	GetCountClassificationPredict() (int, error)
	GetTweetWithoutClassificationPredict(map[string]interface{}) error
	SetTweetClassificationPredict(bson.ObjectId, string, float64)
}
