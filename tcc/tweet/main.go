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

	fmt.Println("Starting tweet with config:")

	fmt.Println("TWITTER_CONSUMER_API_KEY", os.Getenv("TWITTER_CONSUMER_API_KEY"))
	fmt.Println("TWITTER_CONSUMER_API_SECRET", os.Getenv("TWITTER_CONSUMER_API_SECRET"))
	fmt.Println("TWITTER_ACCESS_TOKEN", os.Getenv("TWITTER_ACCESS_TOKEN"))
	fmt.Println("TWITTER_ACCESS_TOKEN_SECRET", os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_API_KEY"), os.Getenv("TWITTER_CONSUMER_API_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
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
	fmt.Println(time.Now(), "Tweet received!")
	s.InsertTweet(message)
	fmt.Println(time.Now(), "Tweet inserted!")
}
