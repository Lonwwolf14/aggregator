package app

import (
	"gator/internal/config"
	"gator/internal/database"
)

type AppState struct {
	AppConfig *config.Config
	DB        *database.Queries
}

func NewAppState(config *config.Config, dbQueries *database.Queries) AppState {
	return AppState{
		AppConfig: config,
		DB:        dbQueries,
	}
}
