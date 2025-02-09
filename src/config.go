package src

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
)

var (
	configDir        string
	configFile       string
	KeyChangerConfig Config
)

func InitConfig() {
	initializeVars()
	reloadConfig()
}

func initializeVars() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configDir = usr.HomeDir + "\\AppData\\Local\\Programs\\key-changer"
	configFile = configDir + "/config.json"

	err = os.MkdirAll(configDir, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	KeyFrom uint32  `json:"key_from"`
	KeyTo   uintptr `json:"key_to"`
}

func reloadConfig() {
	config, err := loadOrCreateConfig()
	if err != nil {
		log.Fatalln(err)
	}

	KeyChangerConfig = config
}

func loadOrCreateConfig() (Config, error) {
	var config Config

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config = Config{
			Keys: []Key{{0x32, 0x33}},
		}

		err := saveConfig(config)
		if err != nil {
			return config, err
		}
		log.Println("File created with default values.")
	} else {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return config, err
		}

		err = json.Unmarshal(data, &config)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}

func saveConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		return err
	}

	log.Println("Config saved.")
	return nil
}
