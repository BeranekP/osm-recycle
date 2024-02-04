package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func ConvertData() {
	log.Println("Converting JSON to GeoJSON")

	path, _ := filepath.Abs("data/containers.json")
	if _, err := os.Stat(path); err != nil {
		log.Printf("File %s not found, fetching", path)
		FetchData()
	}

	sourcePath, err := filepath.Abs("data/containers.json")
	if err != nil {
		log.Fatal(err)
	}

	targetPath, _ := filepath.Abs("data/containers.geojson")
	_, err = os.OpenFile(targetPath, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	source, err := os.ReadFile(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	var r ResponseData

	err = json.Unmarshal(source, &r)
	if err != nil {
		log.Fatal(err)
	}

	geojson := Json2GeoJson(r)
    out, err:=json.MarshalIndent(geojson, "", " ")
    if err!=nil{
    log.Fatal(err)
    }

	os.WriteFile(targetPath, out, 0666)

}
