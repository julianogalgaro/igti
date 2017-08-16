package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/julianogalgaro/classificator/storage"
)

type control struct {
	urlPredictServer string
	storage          storage.Storage
}

var sleepTimeSecond = time.Duration(10)

func (self *control) StartPredict() {

	for {
		err := self.predict()
		if err != nil {
			fmt.Println("Error on predict", err)
			time.Sleep(sleepTimeSecond * time.Second)
		}
	}
}

func (self *control) predict() error {
	tweet := make(map[string]interface{}, 0)
	err := self.storage.GetTweetWithoutClassificationPredict(tweet)
	if err != nil {
		fmt.Println("Error on GetTweetWithoutClassificationPredict", err)
		return err
	}
	tweetJson, err := json.Marshal(tweet)
	fmt.Println("JSON: ", tweet["text"].(string))
	if err != nil {
		fmt.Println("Error on marshal:", err)
		return err
	}
	req, _ := http.NewRequest("POST", self.urlPredictServer, bytes.NewBuffer(tweetJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on predict request", err)
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]interface{}, 0)
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error on parse predict result", err)
		return err
	}

	self.storage.SetTweetClassificationPredict(tweet["_id"].(bson.ObjectId), result["classificationPredict"].(string), result["nnResult"].(float64))

	return nil

}

func NewControl(url string) *control {
	s := storage.NewMongo()
	self := control{
		urlPredictServer: url,
		storage:          s,
	}
	return &self
}
