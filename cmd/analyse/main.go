package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	tweettime "github.com/hareku/tweet-time-analysis"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: [filename]")
	}
	if err := run(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func run(filename string) error {
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

	mp := map[time.Weekday]map[string]int{}
	for _, v := range collection.Tweets {
		t := v.CreatedAt.In(jst).Round(time.Minute * 20)
		if mp[t.Weekday()] == nil {
			mp[t.Weekday()] = make(map[string]int)
		}
		mp[t.Weekday()][t.Format("15:04")]++
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

	for _, wd := range weekdays {
		sum := 0
		keys := make([]string, 0, len(mp[wd]))
		for hm, cnt := range mp[wd] {
			sum += cnt
			keys = append(keys, hm)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		fmt.Printf("%s (%d tweets)\n", wd.String(), sum)

		avg := float64(sum) / float64(len(keys))
		for _, hm := range keys {
			cnt := mp[wd][hm]
			if float64(cnt) < avg {
				continue
			}
			fmt.Printf("[%s] %d\n", hm, cnt)
		}
		fmt.Printf("\n")
	}

	return nil
}
