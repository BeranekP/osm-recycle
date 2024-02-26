package main

import (
	"log"
	"os"
	"time"

	"github.com/BeranekP/osm-recycle/utils"
	"github.com/go-co-op/gocron"
)

func main() {
	args := os.Args
	config := utils.LoadConfig("config/config.json")
    
    // Automatic hourly update
	s := gocron.NewScheduler(time.UTC)
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := s.Every(1).Hour().StartAt(start).Do(update)
	s.StartAsync()
	if err != nil {
		log.Println("Error running scheduled update")
	}

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

func update() {
	log.Println("Running scheduled update")

	config := utils.LoadConfig("config/config.json")
	utils.FetchData(config)
	utils.ConvertData(config)
	utils.ValidateData(config)

	log.Println("Date updated.")

}
