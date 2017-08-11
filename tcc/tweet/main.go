package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/julianogalgaro/tweet/storage"
)

var m storage.Storage

func main() {
	// 	CONSUMER_API_KEY = "CSpU5E2kiJKhoPEIyb58866cp"
	// CONSUMER_API_SECRET = "FGfnCxPHXUuOuxTOVal9Muc6JRHPnUpjAfqBfSCqCbW67jvZ9n"
	// ACCESS_TOKEN = "45275399-OfBm7aZRBavMFjr1GdAsHVBEZhCG9UYxIHUrt8gry"
	// ACCESS_TOKEN_SECRET = "rP1RoTmIyY0LCDmVso8rsh9OrMVrEBTicBSihW3P5fbHA"

	config := oauth1.NewConfig("CSpU5E2kiJKhoPEIyb58866cp", "FGfnCxPHXUuOuxTOVal9Muc6JRHPnUpjAfqBfSCqCbW67jvZ9n")
	token := oauth1.NewToken("45275399-OfBm7aZRBavMFjr1GdAsHVBEZhCG9UYxIHUrt8gry", "rP1RoTmIyY0LCDmVso8rsh9OrMVrEBTicBSihW3P5fbHA")
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	m = storage.NewMongo()

	params := &twitter.StreamFilterParams{
		Track:         []string{"timbrasil", "clarobrasil", "vivoemrede"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		go save(tweet, m)
	}
	demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	stream.Stop()

}

func save(message interface{}, s storage.Storage) {
	s.InsertTweet(message)
	fmt.Println(time.Now(), "Tweet inserted!")
}
