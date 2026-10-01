package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.gatorconfig.json", homeDir), nil
}

func Read() (Config, error) {

	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(configPath)

	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	res, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(res, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
