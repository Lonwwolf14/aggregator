package utils

import (
	"context"
	"encoding/xml"
	"net/http"
)

type RSSFeed struct {
	// Feed struct
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}
	return &feed, nil
}
