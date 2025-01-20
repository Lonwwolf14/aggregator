package app

import "fmt"

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandMap map[string]func(*AppState, Command) error
}

func NewCommands() *Commands {
	return &Commands{CommandMap: make(map[string]func(*AppState, Command) error)}
}

func (c *Commands) Register(name string, handler func(*AppState, Command) error) {
	if name == "" || handler == nil {
		panic("command name or handler cannot be nil")
	}
	c.CommandMap[name] = handler
}

func (c *Commands) Run(s *AppState, cmd Command) error {
	handler, exists := c.CommandMap[cmd.Name]
	if !exists {
		return fmt.Errorf("command %s not found", cmd.Name)
	}
	return handler(s, cmd)
}
