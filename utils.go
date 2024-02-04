package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Update(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found using system")
	}

	fmt.Println("--------------------")
	log.Println("Updating data")
	token := r.URL.Query().Get("token")
	if token != "" {
		systok := os.Getenv("TOKEN")
		valid := token == systok
		if valid {
			FetchData()
			ConvertData()
			ValidateData()
			fmt.Fprint(w, http.StatusOK)
			fmt.Println("--------------------")
			return

		}
	}
	fmt.Fprint(w, http.StatusForbidden)

}

func ServeData() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found using system")
	}
	port := os.Getenv("PORT")

	log.Println("Serving data")
	gjson := http.FileServer(http.Dir("./data"))
	html := http.FileServer(http.Dir("./html"))
	http.Handle("/geojson/", http.StripPrefix("/geojson", gjson))
	http.Handle("/", html)
	http.HandleFunc("/update", Update)
	server := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(server, nil))

}

func ValidateData() {

	var data GeoJson
	source_path, _ := filepath.Abs("data/containers.geojson")
	raw_data, err := os.ReadFile(source_path)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(raw_data, &data)

	if err != nil {
		log.Fatal("Error parsing GeoJSON:", err)
	}

	checked := setValid(&data)
	file, _ := json.MarshalIndent(data, "", " ")
	os.WriteFile(source_path, file, 0644)

	missingRecycling, _ := json.MarshalIndent(checked.missingRecycling, "", " ")
	suspicousTags, _ := json.MarshalIndent(checked.suspiciousTags, "", " ")
	suspicousColor, _ := json.MarshalIndent(checked.suspiciousColor, "", " ")
	withAddress, _ := json.MarshalIndent(checked.withAddress, "", " ")
	fixMe, _ := json.MarshalIndent(checked.fixMe, "", " ")
	stats, _ := json.MarshalIndent(checked.Stats, "", " ")
	missingType, _ := json.MarshalIndent(checked.missingType, "", " ")

	output := OutputData{
		"missingRecycling": missingRecycling,
		"missingType":      missingType,
		"suspicousTags":    suspicousTags,
		"suspicousColor":   suspicousColor,
		"withAddress":      withAddress,
		"fixMe":            fixMe,
		"stats":            stats,
	}
	for k, v := range output {
		partial := fmt.Sprintf("data/%s.geojson", k)
		if k == "stats" {
			partial = fmt.Sprintf("data/%s.json", k)
		}
		filePath, _ := filepath.Abs(partial)
		os.WriteFile(filePath, v, 0644)

	}
}
func FilesExist() bool {
	var validated []string = []string{"containers", "missingType", "missingRecycling", "withAddress", "fixMe", "suspiciousTags", "suspiciousColor"}

	for _, file := range validated {
		path := fmt.Sprintf("data/%s.geojson", file)
		path, _ = filepath.Abs(path)
		if _, err := os.Stat(path); err != nil {
			log.Printf("File %s not found", path)
			return false
		}
	}
	path, _ := filepath.Abs("data/stats.json")
	if _, err := os.Stat(path); err != nil {
		log.Printf("File %s not found", path)
		return false
	}

	log.Println("All data files found")
	return true

}

func Json2GeoJson(r ResponseData) GeoJson {
	output := GeoJson{Type: "FeatureCollection"}
	features := []GeoContainer{}

	for i := range r.Elements {
		container := &r.Elements[i]
		geoContainer := GeoContainer{Type: "Feature"}
		geoContainer.Id = fmt.Sprintf("%s/%d", container.Type, container.Id)
		geoContainer.Geometry.Type = "Point"
		geoContainer.Properties = container.Tags

		if container.Type == "node" {
			geoContainer.Geometry.Coordinates = []float32{container.Lon, container.Lat}
		} else {
			geoContainer.Geometry.Coordinates = []float32{container.Center.Lon, container.Center.Lat}
		}

		features = append(features, geoContainer)
	}
	output.Features = features
	return output

}
