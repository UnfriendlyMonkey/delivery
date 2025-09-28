// Package commands
package commands

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	orderID uuid.UUID
	street string
	volume kernel.Volume

	isValid bool
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume kernel.Volume) (CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return CreateOrderCommand{}, errs.NewValueIsInvalidError("orderID")
	}
	if street == "" {
		return CreateOrderCommand{}, errs.NewValueIsInvalidError("street")
	}
	if !volume.IsValid() {
		return CreateOrderCommand{}, errs.NewValueIsInvalidError("volume")
	}

	return CreateOrderCommand{
		orderID: orderID,
		street: street,
		volume: volume,

		isValid: true,
	}, nil
}

func (c CreateOrderCommand) IsValid() bool {
	return c.isValid
}

func (c CreateOrderCommand) OrderID() uuid.UUID {
	return c.orderID
}

func (c CreateOrderCommand) Street() string {
	return c.street
}

func (c CreateOrderCommand) Volume() kernel.Volume {
	return c.volume
}
