package models

const (
	OrderPending   uint8 = 0
	OrderPayed     uint8 = 1
	OrderReceived  uint8 = 2
	OrderConfirmed uint8 = 3
	OrderRefunding uint8 = 4
	OrderRefunded  uint8 = 5
	OrderCancel    uint8 = 6
)

var statusScope = []uint8{
	OrderPending,
	OrderPayed,
	OrderReceived,
	OrderConfirmed,
	OrderRefunding,
	OrderRefunded,
	OrderCancel,
}

// Order -> Struct to hold information about a specific order from a customer
type Order struct {
	Base
	OrderedBy   *Student   `json:"ordered_by"`
	OrderedFrom *Shop      `json:"ordered_from"`
	OrderItems  []*Product `json:"ordered_items"`
	OrderTotal  float32    `json:"total_price"`
	Status      uint8      `json:"status" gorm:"type:tinyint(1)"`
}
