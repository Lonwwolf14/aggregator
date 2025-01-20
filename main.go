package main

import (
	"database/sql"
	"fmt"
	"gator/internal/app"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/handlers"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	configFile, err := config.Read()
	if err != nil {
		fmt.Printf("Failed to read config: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := sql.Open("postgres", configFile.DbUrl)
	dbQueries := database.New(db)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Initialize AppState and Commands
	appState := app.NewAppState(&configFile, dbQueries)
	commands := app.NewCommands()

	// Register commands
	handlers.RegisterUserHandlers(commands)
	handlers.RegisterRSSHandlers(commands)
	handlers.RegisterFeedFollowHandlers(commands)

	// Parse and execute command-line arguments
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		os.Exit(1)
	}

	command := app.Command{Name: args[0], Args: args[1:]}
	if err := commands.Run(&appState, command); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
