// Package cmd
package cmd

import (
	grpcgeo "delivery/internal/adapters/out/grpc/geo"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/jobs"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs   Config
	gormDB    *gorm.DB
	geoClient ports.GeoClient

	onceGeo sync.Once
	closers []Closer
}

func NewCompositionRoot(configs Config, gormDB *gorm.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
		gormDB:  gormDB,
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

func (cr *CompositionRoot) NewUnitOfWorkFactory() ports.UnitOfWorkFactory {
	factory, err := postgres.NewUnitOfWorkFactory(cr.gormDB)
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}
	return factory
}

func (cr *CompositionRoot) NewCreateOrderHandler() commands.CreateOrderHandler {
	handler, err := commands.NewCreateOrderHandler(cr.NewUnitOfWorkFactory(), cr.NewGeoClient())
	if err != nil {
		log.Fatalf("cannot create CreateOrderHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewCreateCourierHandler() commands.CreateCourierHandler {
	handler, err := commands.NewCreateCourierHandler(cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create CreateCourierHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewAddStoragePlaceHandler() commands.AddStoragePlaceHandler {
	handler, err := commands.NewAddStoragePlaceHandler(cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create NewAddStoragePlaceHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewMoveCouriersHandler() commands.MoveCouriersHandler {
	handler, err := commands.NewMoveCouriersHandler(cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewAssignOrderHandler() commands.AssignOrderHandler {
	handler, err := commands.NewAssignOrderHandler(
		cr.NewUnitOfWorkFactory(), cr.NewOrderDispatcherService(),
	)
	if err != nil {
		log.Fatalf("cannot create AssignOrderHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewGetAllCouriersHandler() queries.GetAllCouriersHandler {
	handler, err := queries.NewGetAllCouriersHandler(cr.gormDB)
	if err != nil {
		log.Fatalf("cannot create GetAllCouriersHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewGetIncompleteOrdersHandler() queries.GetIncompleteOrdersHandler {
	handler, err := queries.NewGetIncompleteOrdersHandler(
		cr.gormDB,
	)
	if err != nil {
		log.Fatalf("cannot create GetIncompleteOrdersHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewAssignOrderJob() cron.Job {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignOrderHandler())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersHandler())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewGeoClient() ports.GeoClient {
	cr.onceGeo.Do(func() {
		client, err := grpcgeo.NewClient(cr.configs.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("cannot create GeoClient: %v", err)
		}
		cr.RegisterCloser(client)
		cr.geoClient = client
	})
	return cr.geoClient
}
