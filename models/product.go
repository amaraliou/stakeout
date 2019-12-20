package models

// Product -> Struct to hold product information
type Product struct {
	Base
	Name       string  `json:"name"`
	BasicPrice float32 `json:"price"`
	SoldBy     *Shop   `json:"sold_by"`
}
