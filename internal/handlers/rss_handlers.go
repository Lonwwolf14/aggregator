package handlers

import (
	"gator/internal/app"
)

func RegisterRSSHandlers(c *app.Commands) {
	c.Register("fetch_feed", handleFetchFeed)
}

func handleFetchFeed(s *app.AppState, cmd app.Command) error {
	// Implementation
	return nil
}
