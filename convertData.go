package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertData() {
	log.Println("Converting JSON to GeoJSON")

	path, _ := filepath.Abs("data/containers.json")
	if _, err := os.Stat(path); err != nil {
		log.Printf("File %s not found, fetching", path)
		FetchData()
	}

	source, err := filepath.Abs("data/containers.json")
	if err != nil {
		log.Fatal(err)
	}

	target, _ := filepath.Abs("data/containers.geojson")
	_, err = os.OpenFile(target, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//defer file.Close()

	//	query := fmt.Sprintf("%s > %s", source, target)

	cmd := exec.Command("osmtogeojson", source)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Error running osmtogeojson", err)
	}
	os.WriteFile(target, out, 0666)

}
