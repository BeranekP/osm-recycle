package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func ConvertData() {
	log.Println("Converting JSON to GeoJSON")

	path, _ := filepath.Abs("data/containers.json.gz")
	if _, err := os.Stat(path); err != nil {
		log.Printf("File %s not found, fetching", path)
		FetchData()
	}

	source := DecompressData("data/containers.json.gz")
	var r ResponseData

	err := json.Unmarshal(source, &r)
	if err != nil {
		log.Fatal("Error unmarshalling: ", err)
	}

	geojson := Json2GeoJson(r)
	out, err := json.Marshal(geojson)
	if err != nil {
		log.Fatal(err)
	}
	CompressData(out, "data/containers.geojson.gz")

	osm, _ := filepath.Abs("data/containers.json.gz")
	os.Remove(osm)

}
