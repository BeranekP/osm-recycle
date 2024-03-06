package types

import "time"

type Config struct {
	Geocode int      `json:"geoCode"`
	Timeout int      `json:"timeOut"`
	Tags    []string `json:"tags"`
	Bad     []string `json:"bad"`
	Common  []string `json:"common"`
	Colors  []string `json:"colors"`
}

type Status struct {
	NoAmenity   bool `json:"noAmenity"`
	NoType      bool `json:"noType"`
	InvalidType bool `json:"invalidType"`
	NoRecycling bool `json:"noRecycling"`
	InvalidTag  bool `json:"invalidTag"`
}

type Container struct {
	Type      string            `json:"type"`
	Id        int               `json:"id"`
	Nodes     []int             `json:"nodes,omitempty"`
	Lat       float32           `json:"lat,omitempty"`
	Lon       float32           `json:"lon,omitempty"`
	Tags      map[string]string `json:"tags"`
	Center    Center            `json:"center,omitempty"`
	User      string            `json:"user,omitempty"`
	Uid       int               `json:"uid,omitempty"`
	Timestamp time.Time         `json:"timestamp,omitempty"`
	Version   int               `json:"version,omitempty"`
}

type Center struct {
	Lon float32 `json:"lon,omitempty"`
	Lat float32 `json:"lat,omitempty"`
}

type TimeStamps struct {
	TimestampOsmBase   string
	TimestampAreasBase string
	Copyright          string
}

type ResponseData struct {
	Version   float32
	Generator string
	Osm3s     TimeStamps
	Elements  []Container `json:"elements"`
}

type GeoJson struct {
	Type     string         `json:"type"`
	Features []GeoContainer `json:"features"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

type GeoContainer struct {
	Type       string            `json:"type"`
	Id         string            `json:"id"`
	Properties map[string]string `json:"properties"`
	Geometry   Geometry          `json:"geometry"`
	Suspicious string            `json:"suspicious"`
	User       string            `json:"user,omitempty"`
	Uid        int               `json:"uid,omitempty"`
	Timestamp  time.Time         `json:"timestamp,omitempty"`
	Recent     bool              `json:"recent,omitempty"`
	Version    int               `json:"version,omitempty"`
}

type CheckedData struct {
	MissingRecycling GeoJson `json:"missingRecycling"`
	MissingType      GeoJson `json:"missingType"`
	MissingAmenity   GeoJson `json:"missingAmenity"`
	SuspiciousTags   GeoJson `json:"suspiciousTags"`
	SuspiciousColor  GeoJson `json:"suspiciousColor"`
	WithAddress      GeoJson `json:"withAddress"`
	FixMe            GeoJson `json:"fixMe"`
	Stats            Stats
	Users            []User
}

type Stats struct {
	Total            int
	MissingRecycling int
	MissingType      int
	MissingAmenity   int
	Fixme            int
	Timestamp        int64
}

type User struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	ValidNew        int    `json:"validNew"`
	ValidModified   int    `json:"validModified"`
	InvalidNew      int    `json:"invalidNew"`
	InvalidModified int    `json:"invalidModified"`
}

type OutputData map[string][]byte
