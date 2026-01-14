package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {

	filePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}

func (cfg *Config) SetUser(userName string) error {

	cfg.CurrentUserName = userName

	if err := write(*cfg); err != nil {
		return err
	}

	return nil

}

func getConfigFilePath() (string, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(filePath, configFileName)

	return fullPath, nil
}

func write(cfg Config) error {

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil

}
