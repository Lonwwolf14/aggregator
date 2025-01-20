package utils

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resFeed RSSFeed
	if err := xml.Unmarshal(body, &resFeed); err != nil {
		return nil, err
	}
	resFeed.Channel.Title = html.UnescapeString(resFeed.Channel.Title)
	resFeed.Channel.Description = html.UnescapeString(resFeed.Channel.Description)

	for i := range resFeed.Channel.Item {
		resFeed.Channel.Item[i].Title = html.UnescapeString(resFeed.Channel.Item[i].Title)
		resFeed.Channel.Item[i].Description = html.UnescapeString(resFeed.Channel.Item[i].Description)
	}
	return &resFeed, nil
}
