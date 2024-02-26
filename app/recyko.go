package main

import (
	"log"
	"os"

	"github.com/BeranekP/osm-recycle/utils"
)

func main() {
	args := os.Args
	log.Println("Loading initial config")
	config := utils.LoadConfig("config/config.json")
	log.Println("Config loaded")
	if len(args) > 1 {
		if args[1] == "-U" {
			utils.FetchData(config)
		}
	}

	if !utils.FilesExist() {
		utils.ConvertData(config)
		utils.ValidateData(config)
	}
	utils.ServeData()

}
