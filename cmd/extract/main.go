package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	c, err := tweettime.NewCollectionFromReader(f)
	if err != nil {
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
	var (
		wd    time.Weekday
		found bool
	)
	for _, v := range weekdays {
		if v.String() == weekdayStr {
			wd = v
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("not found weekday %q", weekdayStr)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		v, err := c.ReadTweet()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("No more tweets.")
				return nil
			}
			return fmt.Errorf("read tweet: %w", err)
		}

		t := v.CreatedAt.In(jst).Round(time.Minute * 20)
		if t.Weekday() == wd && t.Format("15:04") == hm {
			fmt.Print(strings.Repeat("-", 30), "\n")
			fmt.Printf("%s\n", v.Text)
			fmt.Printf("%s (%v days ago)", v.CreatedAt.In(jst).Format(time.RFC3339), math.Round(time.Since(v.CreatedAt).Hours()/24))
			fmt.Printf(" https://twitter.com/%s/status/%s\n\n", c.Meta.UserName, v.ID)

			fmt.Printf("Press any key to continue.\n")
			if !scanner.Scan() {
				break
			}
		}
	}

	return nil
}
