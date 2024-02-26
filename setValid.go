package main

import (
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
	missingAmenity := GeoJson{Type: "FeatureCollection"}
	users := make(map[int]*User)

	var stats Stats
	stats.Timestamp = time.Now().UnixMilli()

	for i := range data.Features {
		container := &data.Features[i]
		stats.Total += 1
		_, ok := users[container.Uid]
		valid := true
		if !ok {
			users[container.Uid] = &User{Name: container.User, Id: container.Uid}
		}

		user, _ := users[container.Uid]

		// check if newer than 1 week
		if container.Timestamp.After(time.Now().Add(time.Duration(-168) * time.Hour)) {
			container.Recent = true
		}

		if (container.Properties["recycling_type"] != "container") && (container.Properties["recycling_type"] != "centre") && (container.Properties["recycling_type"] != "bin") {
			missingType.Features = append(missingType.Features, *container)
			stats.MissingType += 1
			valid = false

		}
		checkSubstance := 0
		for key, value := range container.Properties {
			if strings.Contains(key, "recycling:") && value == "yes" {
				checkSubstance += 1
			}
		}

		if checkSubstance == 0 {

			missingRecycling.Features = append(missingRecycling.Features, *container)
			stats.MissingRecycling += 1
			valid = false

		}
		if container.Properties["amenity"] != "recycling" {
			missingAmenity.Features = append(missingAmenity.Features, *container)
			stats.MissingAmenity += 1
			valid = false
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
				suspiciousTags.Features = append(suspiciousTags.Features, *container)
			}
		}
		if container.Properties["colour"] != "" {
			if !validColor(container.Properties["colour"]) {
				suspiciousColor.Features = append(suspiciousColor.Features, *container)
			}
		}
		if hasAddress(container.Properties) {
			withAddress.Features = append(withAddress.Features, *container)
		}

		if container.Properties["fixme"] != "" {
			fixMe.Features = append(fixMe.Features, *container)
			stats.Fixme += 1

		}

		if valid {
			if container.Version == 1 {
				user.ValidNew += 1
			} else {
				user.ValidModified += 1
			}

		} else {
			if container.Version == 1 {
				user.InvalidNew += 1
			} else {
				user.InvalidModified += 1
			}

		}
	}
	var user_arr []User
	for _, v := range users {
		user_arr = append(user_arr, *v)

	}
	return CheckedData{missingRecycling, missingType, missingAmenity, suspiciousTags, suspiciousColor, withAddress, fixMe, stats, user_arr}
}

func validKeys(props map[string]string, c *GeoContainer) bool {
	valid := []string{"amenity", "recycling_type", "name", "location", "operator",
		"opening_hours", "id", "access", "source", "covered", "wheelchair",
		"description", "check_date:recycling", "ref", "indoor", "collection_times",
		"colour", "check_date", "source:amenity", "website", "note", "material",
		"operator:website", "temporary", "mapillary", "operator:type", "email", "fixme",
		"phone", "mobile", "landuse", "image", "start_date", "fee", "industrial", "level", "count", "charity", "layer"}
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
			strings.HasPrefix(key, "check_date:") ||
			strings.HasPrefix(key, "payment:")) {
			//fmt.Println(key, value)
			c.Suspicious += key
			return false
		}

	}
	return true

}

func validColor(color string) bool {
	valid := []string{"red", "blue", "green", "brown", "yellow", "white", "orange", "gray"}
	colors := strings.Split(color, ";")
	count := 0
	for _, c := range colors {

		if slices.Contains(valid, strings.TrimSpace(c)) {
			count += 1
		}
	}

	return count == len(colors)
}

func hasAddress(props map[string]string) bool {
	for key, _ := range props {
		if strings.HasPrefix(key, "addr:") {
			return true
		}
	}
	return false

}
