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

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUserID.ID,
		FeedID:    dbFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)

	fmt.Println("Feed has been created successfully:")
	if err := printFeed(s, dbFeed); err != nil {
		return err
	}

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

func handlerFetch(s *state, cmd command) error {

	if len(cmd.arguments) != 0 {
		return fmt.Errorf("too many arguments")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerFollow(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("you need to enter an URL to follow a feed")
	}
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("too many arguments to follow a feed")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user information: %w", err)
	}

	currentFeed, err := s.db.GetFeedByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("could not get feed information for following: %w", err)
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    currentFeed.ID,
	}

	feed, err := s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}

	fmt.Printf("%v is now following:\n", feed.UserName)
	fmt.Printf("- %v\n", feed.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user information")
	}

	list, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("could not get following feeds information: %w", err)
	}

	fmt.Printf("You are currently following:\n")

	for _, feed := range list {
		fmt.Printf("- %v\n", feed.FeedName)
	}

	return nil
}

func handlerListAllFeeds(s *state, cmd command) error {

	feeds, err := s.db.ListAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not list any feeds: %w", err)
	}

	for _, feed := range feeds {
		err := printFeed(s, feed)
		if err != nil {
			return err
		}
	}

	return nil
}

func printFeed(s *state, feed database.Feed) error {

	name, err := s.db.GetUserByID(context.Background(), feed.UserID)
	if err != nil {
		return fmt.Errorf("could not get name of user who created feed: %w", err)
	}

	fmt.Printf("* Feed:\n")
	fmt.Printf("* - Name:    %s\n", feed.Name)
	fmt.Printf("* - URL:     %s\n", feed.Url)
	fmt.Printf("* - added by %s\n", name.Name)

	return nil
}
