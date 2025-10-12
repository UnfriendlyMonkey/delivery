package courierrepo

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CourierRepository = &Repository{}

type Repository struct {
	tracker Tracker
}

func NewRepository(tracker Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}

	return &Repository{
		tracker: tracker,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	// r.tracker.Track(aggregate)
	//
	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()
	if tx == nil {
		return errs.NewValueIsRequiredError("transaction not initialized")
	}

	session := &gorm.Session{FullSaveAssociations: true}
	err := tx.WithContext(ctx).Session(session).Create(&dto).Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	// r.tracker.Track(aggregate)

	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()
	if tx == nil {
		return errs.NewValueIsRequiredError("transaction not initialized")
	}

	session := &gorm.Session{FullSaveAssociations: true}
	err := tx.WithContext(ctx).
		Session(session).
		Save(&dto).Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	dto := CourierDTO{}

	tx := r.getTxOrDB()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, ID)
	if result.RowsAffected == 0 {
		return nil, nil
	}

	aggregate := DtoToDomain(dto)

	return aggregate, nil
}

func (r *Repository) GetAllAvailable(ctx context.Context) ([]*courier.Courier, error) {
	var couriers []CourierDTO
	tx := r.getTxOrDB()
	result := tx.WithContext(ctx).
		Preload("StoragePlaces").
		Where("NOT EXISTS (?)",
			tx.Model(&StoragePlaceDTO{}).
				Select("1").
				Where("storage_places.courier_id = couriers.id AND storage_places.order_id IS NOT NULL"),
		).
		Find(&couriers)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("available couriers", nil)
	}
	aggregates := make([]*courier.Courier, len(couriers))
	for i, c := range couriers {
		aggregates[i] = DtoToDomain(c)
	}

	return aggregates, nil
}

func (r *Repository) getTxOrDB() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}
