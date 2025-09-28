package queries

import "github.com/google/uuid"

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}

type LocationResponse struct {
	X int
	Y int
}

func (CourierResponse) TableName() string {
	return "couriers"
}
