package storage

type Storage interface {
	GetTweetWithoutClassification(map[string]interface{}) error
	SetTweetClassification(string, string)
}
