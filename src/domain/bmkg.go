package domain

import (
	"time"
)

// ResponseBmkgAPI represents the structure of the response from the BMKG API.
// ResponseBmkgAPI represents the structure of the response from the BMKG API.
type ResponseBmkgAPI struct {
	Infogempa Infogempa `json:"Infogempa"`
}

// Infogempa represents the information about the earthquake.
type Infogempa struct {
	Gempa Gempa `json:"gempa"`
}

// Gempa represents the details of the earthquake.
type Gempa struct {
	Tanggal     string    `json:"Tanggal"`     // Date of the earthquake
	Jam         string    `json:"Jam"`         // Time of the earthquake
	DateTime    time.Time `json:"DateTime"`    // Date and time of the earthquake
	Coordinates string    `json:"Coordinates"` // Coordinates of the earthquake
	Lintang     string    `json:"Lintang"`     // Latitude of the earthquake
	Bujur       string    `json:"Bujur"`       // Longitude of the earthquake
	Magnitude   string    `json:"Magnitude"`   // Magnitude of the earthquake
	Kedalaman   string    `json:"Kedalaman"`   // Depth of the earthquake
	Wilayah     string    `json:"Wilayah"`     // Region affected by the earthquake
	Potensi     string    `json:"Potensi"`     // Potential impact of the earthquake
	Dirasakan   string    `json:"Dirasakan"`   // Felt impact of the earthquake
	Shakemap    string    `json:"Shakemap"`    // Shakemap of the earthquake
}
