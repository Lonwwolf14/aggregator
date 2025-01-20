package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type AppState struct {
	AppConfig *config.Config
	DB        *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandMap map[string]func(*AppState, Command) error
}

func (c *Commands) register(name string, handler func(*AppState, Command) error) {
	if name == "" {
		panic("command name cannot be empty")
	}
	if handler == nil {
		panic("handler function cannot be nil")
	}
	c.CommandMap[name] = handler
}

func handleRegister(s *AppState, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("username not provided")
	}

	username := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), sql.NullString{String: username, Valid: true})
	if err == nil {
		return fmt.Errorf("user %s already exists", username)
	}

	_, err = s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      sql.NullString{String: username, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create user")
	}
	s.AppConfig.SetUser(username)

	fmt.Printf("User %s added successfully\n", username)

	return nil

}

func handlerLogin(s *AppState, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("username not provided")
	}

	username := cmd.Args[0]
	err := s.AppConfig.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user")
	}
	_, err = s.DB.GetUser(context.Background(), sql.NullString{
		String: username,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to get user")
	}

	fmt.Printf("Logged in as %s\n", username)

	return nil
}

func handleReset(s *AppState, cmd Command) error {
	err := s.DB.DeleteUser(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete user")
	}
	s.AppConfig.CurrentUserName = ""
	fmt.Println("Users reset successfully")
	return nil

}

func handleList(s *AppState, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users")
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

func main() {
	// Load configuration
	configFile, err := config.Read()
	if err != nil {
		fmt.Printf("Failed to read config: %v\n", err)
		os.Exit(1)
	}

	// Connect to the database
	db, err := sql.Open("postgres", configFile.DbUrl)
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize database queries
	dbQueries := database.New(db)

	appState := AppState{
		AppConfig: &configFile,
		DB:        dbQueries,
	}

	// Initialize command map and register commands
	commands := Commands{
		CommandMap: make(map[string]func(*AppState, Command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handleRegister)
	commands.register("reset", handleReset)
	commands.register("users", handleList)

	// Parse command-line arguments
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		os.Exit(1)
	}

	// Execute the command
	command := Command{
		Name: args[0],
		Args: args[1:],
	}
	if err := commands.run(&appState, command); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *Commands) run(s *AppState, cmd Command) error {
	handler, exists := c.CommandMap[cmd.Name]
	if !exists {
		return fmt.Errorf("command %s not found", cmd.Name)
	}

	if err := handler(s, cmd); err != nil {
		return fmt.Errorf("error executing command %s: %w", cmd.Name, err)
	}

	return nil
}
