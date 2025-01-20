package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/app"
	"gator/internal/database"
	"gator/internal/middleware"
	"time"

	"github.com/google/uuid"
)

const (
	ErrorArgsNotFound = "Insufficient Args"
)

func RegisterFeedFollowHandlers(c *app.Commands) {
	c.Register("follow", middlewareLoggedInWrapper(handleFollow))
	c.Register("following", middlewareLoggedInWrapper(handleFollowing))
	c.Register("unfollow", middlewareLoggedInWrapper(handleUnfollow))
}

func middlewareLoggedInWrapper(handler func(*app.AppState, app.Command, database.User) error) func(*app.AppState, app.Command) error {
	return middleware.MiddlewareLoggedIn(handler)
}

func handleFollow(s *app.AppState, cmd app.Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if len(cmd.Args) < 1 {
		println(ErrorArgsNotFound)
		return nil
	}

	feedUrl := cmd.Args[0]
	feed, err := s.DB.GetFeedByURL(ctx, feedUrl)
	if err != nil {
		return err
	}
	feedId := feed.ID
	userId, err := s.DB.GetUserIDByName(ctx, sql.NullString{
		String: s.AppConfig.CurrentUserName,
		Valid:  true,
	})
	if err != nil {
		return err
	}
	_, err = s.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: uuid.NullUUID{UUID: userId, Valid: true},
		FeedID: uuid.NullUUID{UUID: feedId, Valid: true},
	})
	if err != nil {
		return err
	}

	return nil

}

func handleFollowing(s *app.AppState, cmd app.Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId, err := s.DB.GetUserIDByName(ctx, sql.NullString{
		String: s.AppConfig.CurrentUserName,
		Valid:  true,
	})
	if err != nil {
		return err
	}
	follows, err := s.DB.GetFeedFollowsForUser(ctx, uuid.NullUUID{UUID: userId, Valid: true})
	if err != nil {
		return err
	}
	for _, follow := range follows {
		feedId := follow.FeedID.UUID
		feedName, err := s.DB.GetFeedNameById(ctx, feedId)
		if err != nil {
			return err
		}
		fmt.Printf("-- %s\n", feedName)

	}
	return nil

}

func handleUnfollow(s *app.AppState, cmd app.Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if len(cmd.Args) < 1 {
		println(ErrorArgsNotFound)
		return nil
	}
	feedUrl := cmd.Args[0]
	feed, err := s.DB.GetFeedByURL(ctx, feedUrl)
	if err != nil {
		return err
	}
	feedId := feed.ID

	userId, err := s.DB.GetUserIDByName(ctx, sql.NullString{
		String: s.AppConfig.CurrentUserName,
		Valid:  true,
	})
	if err != nil {
		return err
	}
	err = s.DB.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: uuid.NullUUID{UUID: userId, Valid: true},
		FeedID: uuid.NullUUID{UUID: feedId, Valid: true},
	})
	if err != nil {
		fmt.Printf("Error deleting feed follow: %v", err)
		return err
	}
	fmt.Printf("Successfully unfollowed feed %s\n", feedUrl)
	return nil

}
