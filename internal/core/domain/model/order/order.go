// Package order
package order

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

var ErrOrderStatusIsWrongForAction = errors.New("wrong order status for the action")

const (
	VolumeOK = 5
)

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  kernel.Location
	volume    kernel.Volume
	status    Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume kernel.Volume) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsInvalidError("orderID")
	}
	if !location.IsValid() {
		return nil, errs.NewValueIsInvalidError("location")
	}
	if !volume.IsValid() {
		return nil, errs.NewValueIsInvalidError("volume")
	}
	return &Order{
		id:        orderID,
		location:  location,
		volume:    volume,
		status:    StatusCreated,
	}, nil
}

// CreateOrderOK may be used for testing as normal order object w/o errors
func CreateOrderOK() *Order {
	orderID := uuid.New()
	location, _ := kernel.RandomLocation()
	volume, _ := kernel.NewVolume(VolumeOK)
	o, _ := NewOrder(orderID, location, *volume)
	return o
}

func (o *Order) Equal(target *Order) bool {
	if target == nil {
		return false
	}
	return o.id == target.id
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() kernel.Volume {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Assign(courierID *uuid.UUID) error {
	// TODO: do any validation here?
	if courierID == nil || *courierID == uuid.Nil {
		return errs.NewValueIsInvalidError("courierID")
	}
	o.status = StatusAssigned
	o.courierID = courierID
	return nil
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		return ErrOrderStatusIsWrongForAction
	}
	o.status = StatusCompleted
	return nil
}
