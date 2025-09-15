// Package courier is for aggregate Couries and it's entinies
package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

const (
	BagVolume = 10
	BagName = "bag"
)

var ErrStoragePlaceNotEmptyOrLargeEnough = errors.New("storage place is occupied or less than necessary")

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume kernel.Volume
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume kernel.Volume) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if !totalVolume.IsValid() {
		return nil, errs.NewValueIsInvalidError("totalVolume")
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func NewBag() *StoragePlace {
	bag, _ := NewStoragePlace(BagName, BagVolume)
	return bag
}

func (s *StoragePlace) isOccupied() bool {
	return s.orderID != nil
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() kernel.Volume {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) CanStore(volume kernel.Volume) (bool, error) {
	if !volume.IsValid() {
		return false, errs.NewValueIsInvalidError("volume")
	}
	if s.isOccupied() || !volume.FitsTo(&s.totalVolume) {
		return false, nil
	}
	return true, nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume kernel.Volume) error {
	canStore, err := s.CanStore(volume)
	if err != nil {
		return err
	}
	if !canStore {
		return ErrStoragePlaceNotEmptyOrLargeEnough
	}
	s.orderID = &orderID
	return nil
}

func (s *StoragePlace) Clear(orderID uuid.UUID) error {
	if orderID != *s.orderID {
		return errs.NewValueIsInvalidError("orderID")
	}
	s.orderID = nil
	return nil
}

func (s *StoragePlace) Equal(target *StoragePlace) bool {
	if target == nil {
		return false
	}
	return s.id == target.id
}
