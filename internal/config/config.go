package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePAth()
	if err != nil {
		return Config{}, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()

	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func SetUser(userName string) error {
	config, err := Read()
	if err != nil {
		return err
	}
	config.Current_user_name = userName
	return write(config)
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePAth()
	if err != nil {
		return err
	}

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	return json.NewEncoder(configFile).Encode(cfg)
}

func getConfigFilePAth() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return userHomeDir + "/" + configFileName, nil
}
