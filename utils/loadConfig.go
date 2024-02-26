package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/BeranekP/osm-recycle/types"
)

func LoadConfig(path string) types.Config {
	var config types.Config

	log.Println("Loading initial config")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error loading config file: ", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error parsing config file", err)
	}

	log.Println("Config loaded")
	return config
}
