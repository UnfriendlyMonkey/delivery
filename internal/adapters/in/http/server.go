// Package http
package http

import (
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/pkg/errs"
)

type Server struct {
	createOrderHandler commands.CreateOrderHandler
	createCourierHandler commands.CreateCourierHandler
	getAllCouriersHandler queries.GetAllCouriersHandler
	getIncompleteOrdersHandler queries.GetIncompleteOrdersHandler
}

func NewServer(
	createOrderHandler commands.CreateOrderHandler,
	createCourierHandler commands.CreateCourierHandler,
	getAllCouriersHandler queries.GetAllCouriersHandler,
	getIncompleteOrdersHandler queries.GetIncompleteOrdersHandler,
) (*Server, error) {
	if createOrderHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderHandler")
	}
	if createCourierHandler == nil {
		return nil, errs.NewValueIsRequiredError("createCourierHandler")
	}
	if getAllCouriersHandler == nil {
		return nil, errs.NewValueIsRequiredError("getAllCouriersHandler")
	}
	if getIncompleteOrdersHandler == nil {
		return nil, errs.NewValueIsRequiredError("getIncompleteOrdersHandler")
	}

	return &Server{
		createOrderHandler: createOrderHandler,
		createCourierHandler: createCourierHandler,
		getAllCouriersHandler: getAllCouriersHandler,
		getIncompleteOrdersHandler: getIncompleteOrdersHandler,
	}, nil
}
