package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	var tweets []tweettime.Tweet
	if err := json.NewDecoder(f).Decode(&tweets); err != nil {
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
	for _, v := range tweets {
		t := v.CreatedAt.In(jst).Round(time.Minute * 20)
		if t.Weekday() == wd && t.Format("15:04") == hm {
			extracted = append(extracted, v)
		}
	}

	for _, v := range extracted {
		fmt.Print(strings.Repeat("-", 30), "\n")
		fmt.Printf("%s\n", v.Text)
		fmt.Printf("\n\n")
	}

	return nil
}
