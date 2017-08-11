package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julianogalgaro/classificator/storage"
	"github.com/julienschmidt/httprouter"
)

type api struct {
	storage storage.Storage
}

func (self *api) StartServer() {
	router := httprouter.New()
	router.GET("/tweets/classification", self.getTweetsToClassification)
	router.PUT("/tweets/:id", self.setTweetClassification)
	router.ServeFiles("/static/*filepath", http.Dir("./static/"))
	fmt.Println("Listening localhost:80...")
	log.Fatal(http.ListenAndServe(":80", router))
}

func (self *api) getTweetsToClassification(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tweet := make(map[string]interface{}, 0)
	err := self.storage.GetTweetWithoutClassification(tweet)
	if err != nil {
		self.send(w, http.StatusNoContent, nil)
		return
	}

	classified, _ := self.storage.GetCountClassification()
	total, _ := self.storage.GetCountTweets()
	tweet["totalTweets"] = total
	tweet["totalTweetsClassified"] = classified
	self.send(w, http.StatusOK, tweet)
	return
}

func (self *api) send(w http.ResponseWriter, code int, val interface{}) error {
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, TRACE, GET, HEAD, POST, PUT")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(val)
}

func (self *api) setTweetClassification(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	configJson := make(map[string]interface{})
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		self.send(w, http.StatusBadRequest, err)
		return
	}
	if err := r.Body.Close(); err != nil {
		self.send(w, http.StatusInternalServerError, nil)
		return
	}
	strings.Replace(string(body), "\\", "", -1)

	if err := json.Unmarshal(body, &configJson); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("Error encode json")
			self.send(w, http.StatusBadRequest, nil)
			return
		}
	}
	fmt.Println(configJson)
	if configJson["classification"] == nil {
		fmt.Println("Error record without classification")
		self.send(w, http.StatusBadRequest, nil)
		return
	}
	id := ps.ByName("id")
	self.storage.SetTweetClassification(id, configJson["classification"].(string))
	self.send(w, http.StatusOK, nil)

}

func NewApi() *api {
	m := storage.NewMongo()
	self := api{
		storage: m,
	}
	return &self
}
