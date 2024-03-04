package main

import (
	"log"
	"os"
	"time"

	"github.com/BeranekP/osm-recycle/utils"
	"github.com/go-co-op/gocron"
	_ "time/tzdata"
)

func main() {
	args := os.Args
	config := utils.LoadConfig("config/config.json")

	// Automatic hourly update
	s := gocron.NewScheduler(time.UTC)
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	job, err := s.Every(1).Hour().StartAt(start).Tag("update_data").AfterJobRuns(func(jobName string) {
		j, _ := s.FindJobsByTag("update_data")
		log.Println("Next update: ", j[0].NextRun())
	}).Do(update)

	s.StartAsync()
	if s.IsRunning() {
		log.Println("Automatic update scheduled: ", job.NextRun())

	}
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

	log.Println("Data updated")

}
