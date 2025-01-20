package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(homeDir + configFileName)
	if err != nil {
		return Config{}, err
	}
	var configData Config
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return Config{}, err
	}
	return configData, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(homeDir+configFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
