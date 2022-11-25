package tweettime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Tweet struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type TwitterError struct {
	Message string
}

func GetTweets(ctx context.Context, cli *http.Client, uid string, since time.Time, limit int) ([]Tweet, error) {
	var res struct {
		Data []Tweet `json:"data"`
		Meta struct {
			NextToken string `json:"next_token"`
		} `json:"meta"`
		Errors []TwitterError
	}
	tweets := make([]Tweet, 0)
	nextToken := ""
	for rem := limit - len(tweets); rem > 0; rem = limit - len(tweets) {
		q := url.Values{
			"tweet.fields": []string{"text,created_at"},
		}
		if nextToken != "" {
			q.Add("pagination_token", nextToken)
		}

		// max_results must be between 5 and 100.
		max := 100
		if max > rem {
			max = rem
		}
		if max < 5 {
			max = 5
		}
		q.Add("max_results", fmt.Sprint(max))

		req, err := http.NewRequestWithContext(
			ctx, http.MethodGet,
			fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets?%s", uid, q.Encode()),
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("make http request: %w", err)
		}

		resp, err := cli.Do(req)
		if err != nil {
			return nil, fmt.Errorf("http: %w", err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, fmt.Errorf("decode json:%w", err)
		}

		if len(res.Errors) > 0 {
			return nil, fmt.Errorf("twitter error: %s", res.Errors[0].Message)
		}

		for _, v := range res.Data {
			if v.CreatedAt.Before(since) {
				return tweets, nil
			}
			tweets = append(tweets, v)
			if len(tweets) >= limit {
				break
			}
		}

		if res.Meta.NextToken == "" {
			break
		}
		nextToken = res.Meta.NextToken
	}
	return tweets, nil
}
