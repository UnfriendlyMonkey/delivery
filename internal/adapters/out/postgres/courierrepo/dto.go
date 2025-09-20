package courierrepo

import "github.com/google/uuid"

type CourierDTO struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string
	Location      LocationDTO `gorm:"embedded;embeddedPrefix:location_"`
	Speed         int
	StoragePlaces []*StoragePlaceDTO `gorm:"foreignKey:CourierID;constraint:OnDelete:CASCADE;"`
}

type StoragePlaceDTO struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	TotalVolume int
	OrderID     *uuid.UUID `gorm:"type:uuid"`
	CourierID   *uuid.UUID `gorm:"type:uuid;index"`
}

type LocationDTO struct {
	X int
	Y int
}

func (CourierDTO) TableName() string {
	return "couriers"
}

func (StoragePlaceDTO) TableName() string {
	return "storage_places"
}
