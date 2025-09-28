package commands

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/pkg/errs"
)

type CreateCourierCommand struct {
	name  string
	speed int

	isValid bool
}

func NewCreateCourierCommand(name string, speed int) (CreateCourierCommand, error) {
	if name == "" {
		return CreateCourierCommand{}, errs.NewValueIsInvalidError("name")
	}
	if speed < courier.MinSpeed || speed > courier.MaxSpeed {
		return CreateCourierCommand{}, errs.NewValueIsInvalidError("speed")
	}

	return CreateCourierCommand{
		name: name,
		speed: speed,

		isValid: true,
	}, nil
}

func (c CreateCourierCommand) IsValid() bool {
	return c.isValid
}

func (c CreateCourierCommand) Name() string {
	return c.name
}

func (c CreateCourierCommand) Speed() int {
	return c.speed
}
