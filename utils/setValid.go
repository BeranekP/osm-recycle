package utils

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/BeranekP/osm-recycle/types"
)

func setValid(data *types.GeoJson, config types.Config) types.CheckedData {
	log.Println("Validating data")

	missingRecycling := types.GeoJson{Type: "FeatureCollection"}
	suspiciousTags := types.GeoJson{Type: "FeatureCollection"}
	suspiciousColor := types.GeoJson{Type: "FeatureCollection"}
	withAddress := types.GeoJson{Type: "FeatureCollection"}
	fixMe := types.GeoJson{Type: "FeatureCollection"}
	missingType := types.GeoJson{Type: "FeatureCollection"}
	missingAmenity := types.GeoJson{Type: "FeatureCollection"}
	users := make(map[int]*types.User)

	var stats types.Stats
	stats.Timestamp = time.Now().UnixMilli()

	for i := range data.Features {
		container := &data.Features[i]
		stats.Total += 1
		_, ok := users[container.Uid]
		valid := true
		if !ok {
			users[container.Uid] = &types.User{Name: container.User, Id: container.Uid}
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

		if !validKeys(container.Properties, container, config) {
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
			if !validColor(container.Properties["colour"], config) {
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
	var user_arr []types.User
	for _, v := range users {
		user_arr = append(user_arr, *v)

	}
	return types.CheckedData{
		MissingRecycling: missingRecycling,
		MissingType:      missingType,
		MissingAmenity:   missingAmenity,
		SuspiciousTags:   suspiciousTags,
		SuspiciousColor:  suspiciousColor,
		WithAddress:      withAddress,
		FixMe:            fixMe,
		Stats:            stats,
		Users:            user_arr}
}

func validKeys(props map[string]string, c *types.GeoContainer, config types.Config) bool {
	valid := config.Tags

	for key, _ := range props {
		prefixSuffix := strings.Split(key, ":")
		prefix := prefixSuffix[0]
		suffix := ""
		if len(prefixSuffix) > 1 {
			suffix = prefixSuffix[1]

		}
		if slices.Contains(config.Bad, key) {
			c.Suspicious += key
			return false
		}

		if prefix == "recycling" && !slices.Contains(config.Common, suffix) {
			c.Suspicious += fmt.Sprintf("%s:%s", prefix, suffix)
			return false
		}
		if !slices.Contains(valid, prefix) {
			//fmt.Println(key, value)
			c.Suspicious += key
			return false
		}

	}
	return true

}

func validColor(color string, config types.Config) bool {
	valid := config.Colors
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
