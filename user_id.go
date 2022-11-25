package tweettime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUserID(ctx context.Context, cli *http.Client, username string) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		"https://api.twitter.com/2/users/by/username/"+username,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("make http request: %w", err)
	}

	resp, err := cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	var res struct {
		Data struct {
			ID string
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("decode json: %w", err)
	}

	return res.Data.ID, nil
}
