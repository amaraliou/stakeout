package models

// Product -> Struct to hold product information
type Product struct {
	Base
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Code         string  `json:"code"`
	Price        float32 `json:"price"`
	InSale       bool    `json:"is_in_sale"`
	Discount     int     `json:"discount"`
	DiscountUnit string  `json:"discount_unit"`
	SoldBy       *Shop   `json:"sold_by"`
	Reward       int     `json:"reward"`
}
