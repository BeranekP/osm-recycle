package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BeranekP/osm-recycle/types"
	"github.com/joho/godotenv"
)

func Update(w http.ResponseWriter, r *http.Request) {
    log.Println("Loading current config")
    config := LoadConfig("config/config.json")
    log.Println("Config loaded")
	if r.Method == "POST" {

		err := godotenv.Load()
		if err != nil {
			log.Println(".env not found using system")
		}

		fmt.Println("--------------------")
		log.Println("Updating data")
		token := r.FormValue("token")
		if token != "" {
			hasher := sha512.New()
			hasher.Write([]byte(token))
			hashed := hex.EncodeToString(hasher.Sum(nil))

			systok := os.Getenv("TOKEN")
			valid := hashed == systok
			if valid {
				FetchData(config)
				ConvertData(config)
				ValidateData(config)
				fmt.Fprint(w, http.StatusOK)
				fmt.Println("--------------------")
				return

			}
		}
	}
	log.Println("Invalid request")
	fmt.Fprint(w, http.StatusForbidden)

}

func ServeData() {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found using system")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "3333"
	}

	log.Printf("Serving data on http://localhost:%s/", port)
	gjson := http.FileServer(http.Dir("./data"))
	html := http.FileServer(http.Dir("./html"))
	http.Handle("/geojson/", http.StripPrefix("/geojson", gjson))
	http.Handle("/", html)
	http.HandleFunc("/update", Update)
	server := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(server, nil))

}

func ValidateData(config types.Config) {

	var data types.GeoJson
	raw_data := DecompressData("data/containers.geojson.gz")

	err := json.Unmarshal(raw_data, &data)

	if err != nil {
		log.Fatal("Error parsing GeoJSON:", err)
	}

	checked := setValid(&data, config)
	//	file, _ := json.Marshal(data)

	//  CompressData(file,"data/containers.geojson.gz" )

	missingRecycling, _ := json.Marshal(checked.MissingRecycling)
	suspicousTags, _ := json.Marshal(checked.SuspiciousTags)
	suspicousColor, _ := json.Marshal(checked.SuspiciousColor)
	withAddress, _ := json.Marshal(checked.WithAddress)
	fixMe, _ := json.Marshal(checked.FixMe)
	stats, _ := json.Marshal(checked.Stats)
	users, _ := json.Marshal(checked.Users)
	missingType, _ := json.Marshal(checked.MissingType)
	missingAmenity, _ := json.Marshal(checked.MissingAmenity)

	output := types.OutputData{
		"missingRecycling": missingRecycling,
		"missingType":      missingType,
		"missingAmenity":   missingAmenity,
		"suspiciousTags":   suspicousTags,
		"suspiciousColor":  suspicousColor,
		"withAddress":      withAddress,
		"fixMe":            fixMe,
		"stats":            stats,
		"users":            users,
	}
	for k, v := range output {
		partial := fmt.Sprintf("data/%s.geojson.gz", k)
		if k == "stats" || k == "users" {
			partial = fmt.Sprintf("data/%s.json.gz", k)
		}
		CompressData(v, partial)
		//filePath, _ := filepath.Abs(partial)
		//os.WriteFile(filePath, v, 0644)

	}
}
func FilesExist() bool {
	data, _ := filepath.Abs("data")
	if _, err := os.Stat(data); err != nil {
		if err = os.Mkdir("data", 0644); err != nil {
			log.Fatal("Data directory not found, error creating: ", err)
		}
	}

	var validated []string = []string{"containers", "missingType", "missingRecycling", "missingAmenity", "withAddress", "fixMe", "suspiciousTags", "suspiciousColor"}

	for _, file := range validated {
		path := fmt.Sprintf("data/%s.geojson.gz", file)
		path, _ = filepath.Abs(path)
		if _, err := os.Stat(path); err != nil {
			log.Printf("File %s not found", path)
			return false
		}
	}
	path, _ := filepath.Abs("data/stats.json.gz")
	if _, err := os.Stat(path); err != nil {
		log.Printf("File %s not found", path)
		return false
	}

	log.Println("All data files found")
	return true

}

func Json2GeoJson(r types.ResponseData) types.GeoJson {
	output := types.GeoJson{Type: "FeatureCollection"}
	features := []types.GeoContainer{}

	for i := range r.Elements {
		container := &r.Elements[i]
		geoContainer := types.GeoContainer{Type: "Feature"}
		geoContainer.Id = fmt.Sprintf("%s/%d", container.Type, container.Id)
		geoContainer.Geometry.Type = "Point"
		geoContainer.Properties = container.Tags
		geoContainer.User = container.User
		geoContainer.Uid = container.Uid
		geoContainer.Timestamp = container.Timestamp
		geoContainer.Version = container.Version

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

func CompressData(jsonBytes []byte, outPath string) {

	gzPath, _ := filepath.Abs(outPath)
	//_, err := os.OpenFile(gzPath, os.O_CREATE, 0666)
	//if err != nil {
	//		log.Fatal(err)
	//	}

	var fileGZ bytes.Buffer
	zipper := gzip.NewWriter(&fileGZ)

	_, err := zipper.Write(jsonBytes)
	if err != nil {
		log.Fatal("GZIP error: ", err)
	}
	zipper.Close()
	os.WriteFile(gzPath, fileGZ.Bytes(), 0644)

}

func DecompressData(path string) []byte {
	abs_path, _ := filepath.Abs(path)
	log.Println("Decompressing file: ", abs_path)
	zippped, err := os.ReadFile(abs_path)
	if err != nil {
		log.Fatal("ReadFile: ", err)
	}

	unzipped, err := gzip.NewReader(bytes.NewReader(zippped))

	if err != nil {
		log.Fatal("Error decompressing: ", err)
	}
	//var output []byte
	output, err := io.ReadAll(unzipped)
	//_, err = unzipped.Read(output)
	if err != nil {
		log.Fatal("Error reading unzipped: ", err)
	}

	unzipped.Close()
	return output
}
