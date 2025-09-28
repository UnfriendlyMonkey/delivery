package commands

import (
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateOrderHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

type createOrderHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

var _ CreateOrderHandler = &createOrderHandler{}

func NewCreateOrderHandler(
	uowFactory ports.UnitOfWorkFactory,
) (CreateOrderHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsInvalidError("uowFactory")
	}

	return &createOrderHandler{
		uowFactory: uowFactory,
	}, nil
}

func (h *createOrderHandler) Handle(ctx context.Context, command CreateOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	orderAggregate, err := uow.OrderRepository().Get(ctx, command.OrderID())
	if err != nil {
		return err
	}
	if orderAggregate != nil {
		return nil
	}

	location, err := kernel.RandomLocation()
	if err != nil {
		return err
	}

	orderAggregate, err = order.NewOrder(command.OrderID(), location, command.Volume())
	if err != nil {
		return err
	}

	err = uow.OrderRepository().Add(ctx, orderAggregate)
	if err != nil {
		return err
	}

	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
