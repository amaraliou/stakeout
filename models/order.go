package models

// Order -> Struct to hold information about a specific order from a customer
type Order struct {
	Base
	OrderedBy   *Student   `json:"ordered_by"`
	OrderedFrom *Shop      `json:"ordered_from"`
	OrderItems  []*Product `json:"ordered_items"`
	OrderTotal  float32    `json:"total_price"`
	Status      string     `json:"status"`
}
