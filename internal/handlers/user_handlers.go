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

const (
	ErrUserNotProvided = "username not provided"
	ErrUserExists      = "user %s already exists"
	ErrCreateUser      = "failed to create user"
	ErrSetUser         = "failed to set user"
	ErrGetUser         = "failed to get user"
	ErrDeleteUser      = "failed to delete users"
	ErrListUsers       = "failed to get users"
)

func RegisterUserHandlers(c *app.Commands) {
	c.Register("register", handleRegister)
	c.Register("login", handleLogin)
	c.Register("reset", handleReset)
	c.Register("users", handleList)
	c.Register("agg", handleAgg)
}

func validateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	return nil
}

func handleRegister(s *app.AppState, cmd app.Command) error {
	if len(cmd.Args) < 1 {
		println(ErrUserNotProvided)
		return nil
	}

	username := cmd.Args[0]
	if err := validateUsername(username); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.DB.GetUser(ctx, sql.NullString{String: username, Valid: true})
	if err == nil {
		return fmt.Errorf(ErrUserExists, username)
	}

	_, err = s.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      sql.NullString{String: username, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("%s: %v", ErrCreateUser, err)
	}
	s.AppConfig.SetUser(username)

	fmt.Printf("User %s added successfully\n", username)
	return nil
}

func handleLogin(s *app.AppState, cmd app.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Username not provided")
		return nil
	}

	username := cmd.Args[0]
	fmt.Printf("Attempting login for user: %s\n", username)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.AppConfig.SetUser(username)
	if err != nil {
		return fmt.Errorf("%s: %v", ErrSetUser, err)
	}

	_, err = s.DB.GetUser(ctx, sql.NullString{
		String: username,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("%s: %v", ErrGetUser, err)
	}

	fmt.Printf("Logged in as %s\n", username)
	return nil
}

func handleReset(s *app.AppState, cmd app.Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.DB.DeleteUser(ctx)
	if err != nil {
		return fmt.Errorf("%s: %v", ErrDeleteUser, err)
	}
	s.AppConfig.CurrentUserName = ""
	fmt.Println("Users reset successfully")
	return nil
}

func handleList(s *app.AppState, cmd app.Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := s.DB.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("%s: %v", ErrListUsers, err)
	}

	for _, user := range users {
		if user.Name.String == s.AppConfig.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name.String)
			continue
		}
		fmt.Printf("%s\n", user.Name.String)
	}
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
