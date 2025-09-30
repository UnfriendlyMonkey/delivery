package queries

import "github.com/google/uuid"

type GetIncompleteOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}

func (OrderResponse) TableName() string {
	return "orders"
}
