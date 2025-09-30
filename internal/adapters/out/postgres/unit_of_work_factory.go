package postgres

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type unitOfWorkFactory struct {
	db *gorm.DB
}

func NewUnitOfWorkFactory(db *gorm.DB) (ports.UnitOfWorkFactory, error) {
	if db == nil {
		return nil, errs.NewValueIsInvalidError("db")
	}

	return &unitOfWorkFactory{db: db}, nil
}

func (f *unitOfWorkFactory) New(ctx context.Context) (ports.UnitOfWork, error) {
	return NewUnitOfWork(f.db.WithContext(ctx))
}
