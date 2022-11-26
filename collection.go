package tweettime

type Collection struct {
	UserName string  `json:"user_name"`
	UserID   string  `json:"user_id"`
	Tweets   []Tweet `json:"tweets"`
}
