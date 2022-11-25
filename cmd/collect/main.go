package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tweettime "github.com/hareku/tweet-time-analysis"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	usage := "usage: [username] [days] [limit]"
	if len(os.Args) != 4 {
		log.Fatal(usage)
	}
	username := os.Args[1]
	days, err := strconv.ParseInt(os.Args[2], 10, 16)
	if err != nil {
		log.Fatal(usage)
	}
	limit, err := strconv.ParseInt(os.Args[3], 10, 16)
	if err != nil {
		log.Fatal(usage)
	}

	if err := run(username, int(days), int(limit)); err != nil {
		log.Fatal(err)
	}
}

func run(username string, days, limit int) error {
	ctx := context.Background()
	conf := clientcredentials.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	cli := &http.Client{
		Transport: &oauth2.Transport{
			Source: conf.TokenSource(ctx),
			Base:   retryablehttp.NewClient().StandardClient().Transport,
		},
	}

	uid, err := tweettime.GetUserID(ctx, cli, username)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Downloading %d tweets.", limit)
	tweets, err := tweettime.GetTweets(ctx, cli, uid, time.Now().Add(time.Hour*time.Duration(-24)*time.Duration(days)), limit)
	if err != nil {
		log.Fatal(err)
	}
	if len(tweets) == 0 {
		return errors.New("no tweets")
	}

	filename := fmt.Sprintf("out/%s_%dtweets_last%ddays_%s.json", username, len(tweets), days, time.Now().Format("2006-01-02-15-04-05"))
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&tweets); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	log.Printf("Saved %d tweets at %s.", len(tweets), filename)
	return nil
}
