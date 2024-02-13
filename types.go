package main

type Status struct {
	NoAmenity   bool `json:"noAmenity"`
	NoType      bool `json:"noType"`
	InvalidType bool `json:"invalidType"`
	NoRecycling bool `json:"noRecycling"`
	InvalidTag  bool `json:"invalidTag"`
}

type Container struct {
	Type   string            `json:"type"`
	Id     int               `json:"id"`
	Nodes  []int             `json:"nodes,omitempty"`
	Lat    float32           `json:"lat,omitempty"`
	Lon    float32           `json:"lon,omitempty"`
	Tags   map[string]string `json:"tags"`
	Center Center            `json:"center,omitempty"`
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
	Status     Status            `json:"status"`
	Suspicious string            `json:"suspicious"`
}

type CheckedData struct {
	missingRecycling GeoJson
	missingType      GeoJson
	missingAmenity   GeoJson
	suspiciousTags   GeoJson
	suspiciousColor  GeoJson
	withAddress      GeoJson
	fixMe            GeoJson
	Stats            Stats
}

type Stats struct {
	Total            int
	MissingRecycling int
	MissingType      int
	MissingAmenity   int
	Fixme            int
	Timestamp        int64
}

type OutputData map[string][]byte
