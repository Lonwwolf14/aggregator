package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/app"
	"gator/internal/database"
	"time"
)

func MiddlewareLoggedIn(handler func(s *app.AppState, cmd app.Command, user database.User) error) func(*app.AppState, app.Command) error {
	return func(s *app.AppState, cmd app.Command) error {
		if s.AppConfig.CurrentUserName == "" {
			fmt.Println("User not logged in")
			return nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		currentUserName := sql.NullString{
			String: s.AppConfig.CurrentUserName,
			Valid:  true,
		}

		user, err := s.DB.GetUser(ctx, currentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
