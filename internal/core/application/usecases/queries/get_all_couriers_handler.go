// Package queries
package queries

import (
	"context"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type GetAllCouriersHandler interface{
	Handle(context.Context, GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

type getAllCouriersHandler struct {
	db *gorm.DB
}

var _ GetAllCouriersHandler = &getAllCouriersHandler{}

func NewGetAllCouriersHandler(db *gorm.DB) (GetAllCouriersHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsInvalidError("gorm DB")
	}
	return &getAllCouriersHandler{db: db}, nil
}

func (h *getAllCouriersHandler) Handle(ctx context.Context, query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return GetAllCouriersResponse{}, errs.NewValueIsInvalidError("query")
	}

	var couriers []CourierResponse
	res := h.db.Raw("SELECT id, name, location_x, location_y FROM couriers").Scan(&couriers)
	if res.Error != nil {
		return GetAllCouriersResponse{}, res.Error
	}

	return GetAllCouriersResponse{
		Couriers: couriers,
	}, nil
}
