// Package cmd
package cmd

import (
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"log"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs Config
	gormDB  *gorm.DB

	closers []Closer
}

func NewCompositionRoot(configs Config, gormDB *gorm.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
		gormDB: gormDB,
	}
}

func (cr *CompositionRoot) NewOrderDispatcherService() services.OrderDispatcherService {
	orderDispatcherService := services.NewOrderDispatcherService()
	return orderDispatcherService
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.gormDB)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}
