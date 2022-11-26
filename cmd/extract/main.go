package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	tweettime "github.com/hareku/tweet-time-analysis"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("usage: [filename] [weekday] [h:m]")
	}
	if err := run(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		log.Fatal(err)
	}
}

func run(filename, weekdayStr, hm string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	var collection tweettime.Collection
	if err := json.NewDecoder(f).Decode(&collection); err != nil {
		return err
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	weekdays := []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}

	var wd time.Weekday
	for _, v := range weekdays {
		if v.String() == weekdayStr {
			wd = v
			break
		}
	}

	extracted := []tweettime.Tweet{}
	for _, v := range collection.Tweets {
		t := v.CreatedAt.In(jst).Round(time.Minute * 20)
		if t.Weekday() == wd && t.Format("15:04") == hm {
			extracted = append(extracted, v)
		}
	}

	log.Printf("Found %d tweets.", len(extracted))

	scanner := bufio.NewScanner(os.Stdin)
	for i, v := range extracted {
		fmt.Print(strings.Repeat("-", 30), "\n")
		fmt.Printf("%s\n", v.Text)
		fmt.Printf("%s (%v days ago)", v.CreatedAt.In(jst).Format(time.RFC3339), math.Round(time.Since(v.CreatedAt).Hours()/24))
		fmt.Printf(" https://twitter.com/%s/status/%s\n\n", collection.UserName, v.ID)

		if i != len(extracted)-1 {
			fmt.Printf("Press any key to continue.")
			if !scanner.Scan() {
				break
			}
		}
	}

	return nil
}
