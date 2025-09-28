package commands

import (
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierHandler interface{
	Handle(context.Context, CreateCourierCommand) error
}

type createCourierHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewCreateCourierHandler(uowFactory ports.UnitOfWorkFactory) (CreateCourierHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsInvalidError("uowFactory")
	}

	return &createCourierHandler{
		uowFactory: uowFactory,
	}, nil
}

func (h *createCourierHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)


	location, err := kernel.RandomLocation()
	if err != nil {
		return err
	}

	courierAggregate, err := courier.NewCourier(command.Name(), command.Speed(), location)
	if err != nil {
		return err
	}

	err = uow.CourierRepository().Add(ctx, courierAggregate)
	if err != nil {
		return err
	}

	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
