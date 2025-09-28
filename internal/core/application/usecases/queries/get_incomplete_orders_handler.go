package queries

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type GetIncompleteOrdersHandler interface{
	Handle(context.Context, GetIncompleteOrdersQuery) (GetIncompleteOrdersResponse, error)
}

type getIncompleteOrdersHandler struct {
	db *gorm.DB
}

var _ GetIncompleteOrdersHandler = &getIncompleteOrdersHandler{}

func NewGetIncompleteOrdersHandler(db *gorm.DB) (GetIncompleteOrdersHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsInvalidError("gorm DB")
	}

	return &getIncompleteOrdersHandler{db: db}, nil
}

func (h *getIncompleteOrdersHandler) Handle(ctx context.Context, query GetIncompleteOrdersQuery) (GetIncompleteOrdersResponse, error) {
	if !query.IsValid() {
		return GetIncompleteOrdersResponse{}, errs.NewValueIsInvalidError("query")
	}

	var orders []OrderResponse

	res := h.db.Raw(
		"SELECT id, location_x, location_y FROM orders WHERE status IN (?, ?)", order.StatusCreated, order.StatusAssigned,
	).Scan(&orders)
	if res.Error != nil {
		return GetIncompleteOrdersResponse{}, res.Error
	}

	return GetIncompleteOrdersResponse{Orders: orders}, nil
}
