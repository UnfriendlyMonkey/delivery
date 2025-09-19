// Package services
package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
)

var ErrCourierNotFound = errors.New("no suitable couriers for this order")

type OrderDispatcherService interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatcherService = &orderDispatcherService{}

type orderDispatcherService struct {
}

func NewOrderDispatcherService() OrderDispatcherService {
	return &orderDispatcherService{}
}

func (d *orderDispatcherService) Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}
	if len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

	var bestTime float64 = -1
	var winner *courier.Courier

	for _, c := range couriers {
		ok, err := c.CanTakeOrder(order)
		if err!= nil || !ok {
			continue
		}
		time, err := c.CalculateTimeToLocation(order.Location())
		if err != nil {
			continue
		}
		if winner == nil || time < bestTime {
			bestTime = time
			winner = c
		}
	}

	if winner == nil {
		return nil, ErrCourierNotFound
	}

	winnerID := winner.ID()
	err := winner.TakeOrder(order)
	if err != nil {
		return nil, ErrCourierNotFound
	}
	err = order.Assign(&winnerID)
	if err != nil {
		_ = winner.CompleteOrder(order)
		return nil, ErrCourierNotFound
	}

	return winner, nil
}
