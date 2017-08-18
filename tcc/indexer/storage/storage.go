package storage

type Storage interface {
	GetTweetsWithClassificationPredictByDate(int) ([]map[string]interface{}, error)
	SetLastProcessingDate(string, int64)
	GetLastProcessingDate(string) (int64, error)
}
