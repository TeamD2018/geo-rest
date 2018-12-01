package models

type OSMEntity struct {
	OSMID   int
	OSMType string
}

type Coordinates [][2]float64

type Geojson struct {
	Coordinates []Coordinates
}

type NominatimReverseResponse struct {
	Geojson *Geojson `json:"geojson"`
}

type LookupResp struct {
	DisplayName string `json:"display_name"`
}

type NominatimLookupResponse []*LookupResp
