package main

import (
	"context"
	"encoding/xml"
	"fmt"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	readableData := RSSFeed{}
	err = xml.Unmarshal(data, &readableData)
	if err != nil {
		return nil, fmt.Errorf("error unmarhsaling code: %w", err)
	}

	readableData.Channel.Title = html.UnescapeString(readableData.Channel.Title)
	readableData.Channel.Description = html.UnescapeString(readableData.Channel.Description)

	for i := range readableData.Channel.Item {
		readableData.Channel.Item[i].Title = html.UnescapeString(readableData.Channel.Item[i].Title)
		readableData.Channel.Item[i].Description = html.UnescapeString(readableData.Channel.Item[i].Description)
	}

	return &readableData, nil

}
