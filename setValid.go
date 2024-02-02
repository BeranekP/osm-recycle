package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"
)

func setValid(data *GeoJson) CheckedData {
	log.Println("Validating data")

	missingRecycling := GeoJson{Type: "FeatureCollection"}
	suspiciousTags := GeoJson{Type: "FeatureCollection"}
	suspiciousColor := GeoJson{Type: "FeatureCollection"}
	withAddress := GeoJson{Type: "FeatureCollection"}
	fixMe := GeoJson{Type: "FeatureCollection"}
	missingType := GeoJson{Type: "FeatureCollection"}
	var stats Stats
	stats.Timestamp = time.Now().UnixMilli()

	for i := range data.Features {
		container := &data.Features[i]
		stats.Total += 1
		if container.Properties["recycling_type"] == "" {
			fmt.Println("MissingRecyclingType")
			container.Status.InvalidType = true
			missingType.Features = append(missingType.Features, *container)

		}
		if container.Properties["amenity"] == "" {
			container.Status.NoAmenity = true
		}
		if (container.Properties["recycling_type"] != "container") && (container.Properties["recycling_type"] != "centre") || (container.Properties["recycling_type"] != "bin") {
			container.Status.InvalidType = true

		}
		checkSubstance := 0
		for key, value := range container.Properties {
			if strings.Contains(key, "recycling:") && value == "yes" {
				checkSubstance += 1
			}
		}

		if checkSubstance == 0 {

			container.Status.NoRecycling = true
			missingRecycling.Features = append(missingRecycling.Features, *container)
			stats.MissingRecycling += 1

		}

		if !validKeys(container.Properties, container) {
			if container.Properties["recycling_type"] == "centre" && (container.Properties["barrier"] == "fence" || container.Properties["barrier"] == "wall") {
				if container.Suspicious == "fence" {
					container.Suspicious = ""
				}
			} else if container.Properties["recycling_type"] == "centre" && container.Properties["building"] != "" {
				if container.Suspicious == "building" {
					container.Suspicious = ""
				}
			} else {
				container.Status.InvalidTag = true
				suspiciousTags.Features = append(suspiciousTags.Features, *container)
			}
		}
		if container.Properties["colour"] != "" {
			if !validColor(container.Properties["colour"]) {
				container.Status.InvalidTag = true
				suspiciousColor.Features = append(suspiciousColor.Features, *container)
			}
		}
		if hasAddress(container.Properties) {
			withAddress.Features = append(withAddress.Features, *container)
		}

		if container.Properties["fixme"] != "" {
			fixMe.Features = append(fixMe.Features, *container)

		}

	}
	return CheckedData{missingRecycling, missingType, suspiciousTags, suspiciousColor, withAddress, fixMe, stats}
}

func validKeys(props map[string]string, c *GeoContainer) bool {
	valid := []string{"amenity", "recycling_type", "name", "location", "operator",
		"opening_hours", "id", "access", "source", "covered", "wheelchair",
		"description", "check_date:recycling", "ref", "indoor", "collection_times",
		"colour", "check_date", "source:amenity", "website", "note", "material",
		"operator:website", "temporary", "mapillary", "operator:type", "email", "fixme",
        "phone", "mobile", "landuse", "image", "start_date"}
	for key, _ := range props {
		if !(strings.HasPrefix(key, "recycling:") ||
			slices.Contains(valid, key) ||
			strings.HasPrefix(key, "collection_times:") ||
			strings.HasPrefix(key, "ref:") ||
			strings.HasPrefix(key, "description:") ||
			strings.HasPrefix(key, "addr:") ||
			strings.HasPrefix(key, "name:") ||
			strings.HasPrefix(key, "contact:") ||
			strings.HasPrefix(key, "ipr:") ||
			strings.HasPrefix(key, "survey:") ||
			strings.HasPrefix(key, "source:") ||
			strings.HasPrefix(key, "check_date:")) {
			//fmt.Println(key, value)
			c.Suspicious += key
			return false
		}

	}
	return true

}

func validColor(color string) bool {
	valid := []string{"red", "blue", "green", "brown", "yellow", "white", "orange", "gray"}
	if slices.Contains(valid, color) {
		return true

	}

	return false
}

func hasAddress(props map[string]string) bool {
	for key, _ := range props {
		if strings.HasPrefix(key, "addr:") {
			return true
		}
	}
	return false

}
