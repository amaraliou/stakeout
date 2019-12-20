package models

// Shop -> Struct to hold shop information (try to figure out how to handle shop reg.)
type Shop struct {
	Base
	Name    string   `json:"name"`
	Address *Address `json:"address"`
}
