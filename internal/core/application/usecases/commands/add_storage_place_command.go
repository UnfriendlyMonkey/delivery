package commands

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type AddStoragePlaceCommand struct {
	courierID uuid.UUID
	name      string
	volume    kernel.Volume

	isValid bool
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, volume kernel.Volume) (AddStoragePlaceCommand, error) {
	if courierID == uuid.Nil {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("courierID")
	}
	if name == "" {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("name")
	}
	if !volume.IsValid() {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("volume")
	}

	return AddStoragePlaceCommand{
		courierID: courierID,
		name: name,
		volume: volume,

		isValid: true,
	}, nil
}

func (c AddStoragePlaceCommand) IsValid() bool {
	return c.isValid
}

func (c AddStoragePlaceCommand) CourierID() uuid.UUID {
	return c.courierID
}

func (c AddStoragePlaceCommand) Name() string {
	return c.name
}

func (c AddStoragePlaceCommand) Volume() kernel.Volume {
	return c.volume
}
