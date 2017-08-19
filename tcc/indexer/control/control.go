package control

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/julianogalgaro/indexer/mappings"
	"github.com/julianogalgaro/indexer/storage"
)

type control struct {
	storage          storage.Storage
	elasticSearchUrl string
	indexName        string
	fieldsToIndexing []string
}

var (
	sleepTimeSecondError   = time.Duration(10)
	sleepTimeSecondSuccess = time.Duration(2)
	limitReturnRecords     = 1000
)

func (self *control) StartIndexer() {
	fmt.Printf("Create index [%s] on [%s]...\n", self.indexName, self.elasticSearchUrl)
	err := self.createElasticIndex()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Index [%s] [OK]\n", self.indexName)

	fmt.Printf("Starting indexing to [%s] on [%s]...\n", self.indexName, self.elasticSearchUrl)
	for {
		err := self.indexing()
		if err != nil {
			fmt.Println("Error on indexing", err)
			time.Sleep(sleepTimeSecondError * time.Second)
		} else {
			time.Sleep(sleepTimeSecondSuccess * time.Second)
		}
	}
}

func (self *control) indexing() error {
	tweets, err := self.storage.GetTweetsWithClassificationPredictByDate(limitReturnRecords)
	if err != nil {
		fmt.Println("Error on GetTweetsWithClassificationPredictByDate", err)
		return err
	}

	for _, tweet := range tweets {
		self.treatTweetToIndexing(tweet)

		tweetJson, err := json.Marshal(tweet)
		if err != nil {
			fmt.Println("Error on marshal:", err)
			return err
		}

		err = self.sendDataToElastic(tweetJson, "tweet", tweet["idstr"].(string))
		if err != nil {
			return err
		}
		self.storage.SetLastProcessingDate("tweet", tweet["classificationPredictDate"].(int64))
	}

	return nil
}

func (self *control) treatTweetToIndexing(t map[string]interface{}) {
	for k, _ := range t {
		i := sort.SearchStrings(self.fieldsToIndexing, k)
		if i < len(self.fieldsToIndexing) && self.fieldsToIndexing[i] == k {
			continue
		} else {
			delete(t, k)
		}
	}

	t["classificationDateStr"] = unixNanoToStringDate(t["classificationDate"], "2006-01-02 15:04:05")
	t["classificationPredictDateStr"] = unixNanoToStringDate(t["classificationPredictDate"], "2006-01-02 15:04:05")

	if t["classificationDateStr"].(string) == "" {
		delete(t, "classificationDateStr")
	}
}

func unixNanoToStringDate(u interface{}, format string) string {
	if u != nil {
		i := u.(int64)
		if i < 16000000000 {
			i = i * 1000000000
		}
		return time.Unix(0, i).Format(format)
	}
	return ""
}

func (self *control) sendDataToElastic(data []byte, typeIndex, id string) error {
	url := self.elasticSearchUrl + self.indexName + "/" + typeIndex + "/" + id

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]interface{}, 0)
	json.Unmarshal(body, &result)

	if result["result"] != nil && (result["result"].(string) == "created" || result["result"].(string) == "updated") {
		return nil
	} else {
		return errors.New("Indexing error" + string(body) + " json: " + string(data))
	}
}

func (self *control) createElasticIndex() error {
	url := self.elasticSearchUrl + self.indexName
	mapping, err := mappings.GetIndexMapping(self.indexName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(mapping))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]interface{}, 0)
	json.Unmarshal(body, &result)

	if (result["acknowledged"] != nil && result["acknowledged"].(bool) == true) ||
		(result["error"] != nil && result["error"].(map[string]interface{})["type"].(string) == "index_already_exists_exception") {
		return nil
	} else {
		return errors.New("Create index error response: " + string(body) + " json: " + string(mapping))
	}
}

func NewControl() *control {
	elasticUrl := "http://localhost:9200/"
	indexName := "twitter"

	s := storage.NewMongo()

	fieldsToIndexing := []string{"idstr", "lang", "retweetcount", "retweeted", "createdat",
		"coordinates", "text", "classification", "classificationDate", "classificationPredictDate",
		"classificationPredict", "classificationPredictRate", "user", "entities"}
	sort.Strings(fieldsToIndexing)

	self := control{
		storage:          s,
		elasticSearchUrl: elasticUrl,
		indexName:        indexName,
		fieldsToIndexing: fieldsToIndexing,
	}
	return &self
}
