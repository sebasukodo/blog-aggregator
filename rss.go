package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sebasukodo/blog-aggregator/internal/database"
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

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.arguments) != 2 {
		return fmt.Errorf("not enough arguments to add feed")
	}
	if !strings.HasPrefix(cmd.arguments[1], "https://") {
		return fmt.Errorf("second argument needs to be a valid url")
	}

	currentUserID, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    currentUserID.ID,
	}

	dbFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("could not add feed: %w", err)
	}

	fmt.Println("Feed has been created successfully:")
	printFeed(dbFeed)

	return nil
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

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:      %s\n", feed.ID)
	fmt.Printf("* Name:    %s\n", feed.Name)
	fmt.Printf("* URL:     %s\n", feed.Url)
	fmt.Printf("* User ID: %s\n", feed.UserID)
}
