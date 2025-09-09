// Package courier is for aggregate Couries and it's entinies
package courier

import (
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

var ErrStoragePlaceNotEmptyOrLargeEnough = errors.New("storage place is occupied or less than necessary")

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if totalVolume <= 0 {
		return nil, errs.NewValueIsInvalidError("totalVolume")
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}, nil
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

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if volume <= 0 {
		return false, errs.NewValueIsInvalidError("volume")
	}
	if s.isOccupied() || s.totalVolume < volume {
		return false, nil
	}
	return true, nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int) error {
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
