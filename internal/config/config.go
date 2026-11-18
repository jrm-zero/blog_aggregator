package config

import (
	"os"
	"encoding/json"
	//"fmt"
	//"strings"
	//"io"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c Config) SetUser(user string) error {
	c.Current_user_name = user
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0666); err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	r, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	
	var config Config
	if err := json.Unmarshal(r, &config); err != nil {
		return Config{}, err
	}
	
	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := homeDir + "/" + configFileName

	return configFilePath, nil
}

//func write(cfg Config) error {
//
//}

