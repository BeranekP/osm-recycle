package utils

import (
	"encoding/json"
	"fmt"
	"github.com/BeranekP/osm-recycle/types"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func FetchData(config types.Config) {
	client := http.Client{}
	timeout := config.Timeout
	id := config.Geocode
	query := fmt.Sprintf(`[out:json][timeout:%d];
                area(id:%d)->.searchArea;
                (nwr[~"^recycling:.*$"~"."](area.searchArea);
                nwr["amenity"="recycling"](area.searchArea);
                nwr["recycling_type"~".*"](area.searchArea);); 
                out center meta;`, timeout, id)

	form := url.Values{}
	form.Add("data", query)
	body := form.Encode()

	req, err := http.NewRequest("POST", "https://overpass-api.de/api/interpreter/", strings.NewReader(body))

	if err != nil {
		log.Fatal("Error creating request", err)
	}
	log.Println("Making request to overpass-api")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error making request", err)
	}

	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error in response", err)
	}

	//fmt.Println(string(payload))
	var containers types.ResponseData
	err = json.Unmarshal(payload, &containers)
	if err != nil {
		log.Fatal("Error parsing data:", err)
	}
	//fmt.Println(containers)
	log.Println("Saving data")
	file, _ := json.Marshal(containers)
	//target, _ := filepath.Abs("data/containers.json")
	CompressData(file, "data/containers.json.gz")
	//err = os.WriteFile(target, file, 0644)
	//if err != nil {
	//	log.Fatal("Error writing file: ", err)
	//}

}
