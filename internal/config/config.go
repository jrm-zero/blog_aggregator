package config

import os

const configFileName = ".gatorconfig.json"

//type Config struct {
//	db_url string `json:"db_url`
//}

func Read() (string, error) {
	filePath := getConfigFilePath()
	return filePath, nil
}

//func SetUser() {
//
//}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configFilePath := homeDir + "\\" + homeDir

	return configFilePath
}

//func write(cfg Config) error {
//
//}

