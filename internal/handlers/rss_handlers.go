package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/app"
	"gator/internal/database"
	"gator/internal/utils"
	"time"

	"github.com/google/uuid"
)

func RegisterRSSHandlers(c *app.Commands) {
	c.Register("fetch_feed", handleFetchFeed)
	c.Register("agg", handleAgg)
	c.Register("addfeed", handleAddFeed)
	c.Register("feeds", handleListFeeds)
}

func handleFetchFeed(s *app.AppState, cmd app.Command) error {
	// Implementation
	return nil
}

func handleAgg(s *app.AppState, cmd app.Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rssFeed *utils.RSSFeed
	rssFeed, err := utils.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}
	fmt.Println("Title:", rssFeed.Channel.Title)
	fmt.Println("Link:", rssFeed.Channel.Link)
	fmt.Println("Description:", rssFeed.Channel.Description)
	fmt.Println("Items:")
	fmt.Printf("Title: %s\n", rssFeed.Channel.Item[0].Title)
	fmt.Printf("Link: %s\n", rssFeed.Channel.Item[0].Link)
	fmt.Printf("Description: %s\n", rssFeed.Channel.Item[0].Description)

	return nil

}

func handleAddFeed(s *app.AppState, cmd app.Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(cmd.Args) < 2 {
		fmt.Printf("Usage: %s addfeed <name> <url>\n", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	user, err := s.DB.GetUser(ctx, sql.NullString{
		String: s.AppConfig.CurrentUserName,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	_, err = s.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to add feed: %v", err)
	}
	fmt.Printf("User %s added Feed %s successfully\n", s.AppConfig.CurrentUserName, name)
	return nil
}

func handleListFeeds(s *app.AppState, cmd app.Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	feeds, err := s.DB.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}

	for _, feed := range feeds {
		user, err := s.DB.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user: %v", err)
		}
		fmt.Printf("Feed: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", user.Name.String)
		fmt.Println("---")
	}
	return nil

}
