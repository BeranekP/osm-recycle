package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func FetchData() {
	client := http.Client{}
	timeout := 60
	geocodes := getGeocodes()
	id := geocodes.CZ
	query := fmt.Sprintf(`[out:json][timeout:%d];
                area(id:%d)->.searchArea;
                nwr["amenity" = "recycling"](area.searchArea);
                out center;`, timeout, id)
	query = fmt.Sprintf(`[out:json][timeout:%d];
                area(id:%d)->.searchArea;
                (nwr[~"^recycling:.*$"~"."](area.searchArea);
                nwr["amenity"="recycling"](area.searchArea);); 
                out center;`, timeout, id)

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
	var containers ResponseData
	err = json.Unmarshal(payload, &containers)
	if err != nil {
		log.Fatal("Error parsing data:", err)
	}
	//fmt.Println(containers)
	log.Println("Saving data")
	file, _ := json.Marshal(containers)
	target, _ := filepath.Abs("data/containers.json")
	err = os.WriteFile(target, file, 0644)
	if err != nil {
		log.Fatal("Error writing file: ", err)
	}

}
