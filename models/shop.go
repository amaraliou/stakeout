package models

// Shop -> Struct to hold shop information (try to figure out how to handle shop reg.)
type Shop struct {
	Base
	Name      string   `json:"name"`
	AddressID uint     `json:"address_id"`
	Address   *Address `json:"address"`
	Admins    []*Admin `json:"admins"`
}
