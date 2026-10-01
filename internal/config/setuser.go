package config

import (
	"encoding/json"
	"os"
)

func write(config Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SetUser(conf Config, current_user_name string) error {
	conf.CurrentUserName = current_user_name

	err := write(conf)
	if err != nil {
		return err
	}

	return nil
}
