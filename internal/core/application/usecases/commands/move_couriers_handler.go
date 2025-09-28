package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type MoveCouriersHandler interface {
	Handle(context.Context, MoveCouriersCommand) error
}

type moveCouriersHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

var _ MoveCouriersHandler = &moveCouriersHandler{}

func NewMoveCouriersHandler(uowFactory ports.UnitOfWorkFactory) (MoveCouriersHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsInvalidError("uowFactory")
	}

	return &moveCouriersHandler{
		uowFactory: uowFactory,
	}, nil
}

func (h *moveCouriersHandler) Handle(ctx context.Context, command MoveCouriersCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsInvalidError("command")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.RollbackUnlessCommitted(ctx)

	orders, err := uow.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		courier, err := uow.CourierRepository().Get(ctx, *order.CourierID())
		if err != nil {
			return err
		}
		err = courier.Move(order.Location())
		if err != nil {
			return err
		}
		if courier.Location().Equal(order.Location()) {
			err = order.Complete()
			if err != nil {
				return err
			}
			err = courier.CompleteOrder(order)
			if err != nil {
				return err
			}
		}
		err = uow.OrderRepository().Update(ctx, order)
		if err != nil {
			return err
		}
		err = uow.CourierRepository().Update(ctx, courier)
		if err != nil {
			return err
		}
	}

	err = uow.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
