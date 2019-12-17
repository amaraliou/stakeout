package models

// CurrentLocation -> struct to hold current location of a student
type CurrentLocation struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
