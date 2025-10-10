package commands

import (
	"context"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AssignOrderHandler interface {
	Handle(context.Context, AssignOrderCommand) error
}

type assignOrderHandler struct {
	uowFactory ports.UnitOfWorkFactory
	dispatcher services.OrderDispatcherService
}

var _ AssignOrderHandler = &assignOrderHandler{}

func NewAssignOrderHandler(
	uowFactory ports.UnitOfWorkFactory, dispatcher services.OrderDispatcherService,
) (AssignOrderHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsInvalidError("uowFactory")
	}
	if dispatcher == nil {
		return nil, errs.NewValueIsInvalidError("dispatcher")
	}

	return &assignOrderHandler{
		uowFactory: uowFactory,
		dispatcher: dispatcher,
	}, nil
}

func (h *assignOrderHandler) Handle(ctx context.Context, command AssignOrderCommand) error {
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

	order, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		return err
	}
	availableCouriers, err := uow.CourierRepository().GetAllAvailable(ctx)
	if err != nil {
		return err
	}
	assignedCourier, err := h.dispatcher.Dispatch(order, availableCouriers)
	if err != nil {
		return err
	}
	err = uow.OrderRepository().Update(ctx, order)
	if err != nil {
		return err
	}
	err = uow.CourierRepository().Update(ctx, assignedCourier)
	if err != nil {
		return err
	}
	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
