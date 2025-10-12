package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AddStoragePlaceHandler interface {
	Handle(context.Context, AddStoragePlaceCommand) error
}

type addStoragePlaceHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

var _ AddStoragePlaceHandler = &addStoragePlaceHandler{}

func NewAddStoragePlaceHandler(uowFactory ports.UnitOfWorkFactory) (AddStoragePlaceHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsInvalidError("uowFactory")
	}

	return &addStoragePlaceHandler{
		uowFactory: uowFactory,
	}, nil
}

func (h *addStoragePlaceHandler) Handle(ctx context.Context, command AddStoragePlaceCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	// Start transaction
	uow.Begin(ctx)

	courierAggregate, err := uow.CourierRepository().Get(ctx, command.CourierID())
	if err != nil {
		return err
	}
	if courierAggregate == nil {
		return errs.NewObjectNotFoundError("courier", command.CourierID())
	}

	err = courierAggregate.AddStoragePlace(command.Name(), int(command.Volume()))
	if err != nil {
		return err
	}

	err = uow.CourierRepository().Update(ctx, courierAggregate)
	if err != nil {
		return err
	}

	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
