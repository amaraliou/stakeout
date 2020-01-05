package models

// Shop -> Struct to hold shop information (try to figure out how to handle shop reg.)
type Shop struct {
	Base
	Name        string   `json:"name"`
	Logo        string   `json:"logo_link"` // Amazon S3
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Description string   `json:"description"`
	AddressID   uint     `json:"address_id"`
	Address     *Address `json:"address"`
}
