package tweettime

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Collection struct {
	Meta CollectionMeta `json:"meta"`
	s    *bufio.Scanner
}

type CollectionMeta struct {
	UserName string `json:"user_name"`
	UserID   string `json:"user_id"`
}

func NewCollectionFromReader(r io.Reader) (*Collection, error) {
	c := &Collection{
		s: bufio.NewScanner(r),
	}
	if !c.s.Scan() {
		return nil, fmt.Errorf("no data from reader")
	}
	if err := json.NewDecoder(bytes.NewReader(c.s.Bytes())).Decode(&c.Meta); err != nil {
		return nil, fmt.Errorf("decode meta: %w", err)
	}
	return c, nil
}

func (c *Collection) ReadTweet() (*Tweet, error) {
	if !c.s.Scan() {
		return nil, io.EOF
	}
	var t Tweet
	if err := json.NewDecoder(bytes.NewReader(c.s.Bytes())).Decode(&t); err != nil {
		return nil, fmt.Errorf("decode tweet: %w", err)
	}
	return &t, nil
}
