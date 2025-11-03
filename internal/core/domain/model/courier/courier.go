package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"math"

	"github.com/google/uuid"
)

var (
	ErrCourierCanNotTakeOrder = errors.New("this Courier can not take this order")
	ErrNoOrderFound           = errors.New("this Courier doesn't carry such order")
)

const (
	MinSpeed = 1
	MaxSpeed = 5
	SpeedOK  = 2
	NameOK   = "SomeName"
)

type Courier struct {
	id            uuid.UUID
	name          string
	location      kernel.Location
	speed         int
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsInvalidError("name")
	}
	if !location.IsValid() {
		return nil, errs.NewValueIsInvalidError("location")
	}
	if speed < MinSpeed || speed > MaxSpeed {
		return nil, errs.NewValueIsOutOfRangeError("speed", speed, MinSpeed, MaxSpeed)
	}
	return &Courier{
		id:       uuid.New(),
		name:     name,
		speed:    speed,
		location: location,
		storagePlaces: []*StoragePlace{
			NewBag(),
		},
	}, nil
}

// RestoreCourier creates from DB record, so no error is expected here
func RestoreCourier(name string, speed int, location kernel.Location, id uuid.UUID, storagePlaces []*StoragePlace) *Courier {
	return &Courier{
		id:            id,
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: storagePlaces,
	}
}

// CreateCourierOK may be used for test as normal courier object w/o errors
func CreateCourierOK() *Courier {
	location, _ := kernel.RandomLocation()
	c, _ := NewCourier(NameOK, SpeedOK, location)
	return c
}

func (c *Courier) Equal(target *Courier) bool {
	if target == nil {
		return false
	}
	return c.id == target.id
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) StoragePlaces() []*StoragePlace {
	return c.storagePlaces
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	v, err := kernel.NewVolume(volume)
	if err != nil {
		return err
	}
	sp, err := NewStoragePlace(name, *v)
	if err != nil {
		return err
	}
	c.storagePlaces = append(c.storagePlaces, sp)
	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if order == nil {
		return false, errs.NewValueIsInvalidError("order")
	}
	for _, place := range c.storagePlaces {
		canStore, err := place.CanStore(order.Volume())
		if err != nil {
			return false, err
		}
		if canStore {
			return true, nil
		}
	}
	return false, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	canTake, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}
	if !canTake {
		return ErrCourierCanNotTakeOrder
	}
	for _, place := range c.storagePlaces {
		canStore, _ := place.CanStore(order.Volume())
		if canStore {
			err = place.Store(order.ID(), order.Volume())
			if err != nil {
				return err
			}
			courierID := c.id
			err = order.Assign(&courierID)
			if err != nil {
				_ = place.Clear(order.ID())
				return err
			}
			return nil
		}
	}
	return ErrCourierCanNotTakeOrder
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	sp, err := c.findStoragePlaceByOrderID(order.ID())
	if err != nil {
		return err
	}
	err = sp.Clear(order.ID())
	if err != nil {
		return err
	}
	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if !target.IsValid() {
		return 0, errs.NewValueIsInvalidError("location")
	}
	dist, err := c.location.Distance(target)
	if err != nil {
		return 0, err
	}
	time := math.Ceil(float64(dist) / float64(c.speed))
	return time, nil
}

func (c *Courier) Move(target kernel.Location) error {
	if !target.IsValid() {
		return errs.NewValueIsInvalidError("location")
	}
	dx := float64(target.X()) - float64(c.location.X())
	dy := float64(target.Y()) - float64(c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	// Calculate new coordinates using int to handle negative values correctly
	newXInt := int(c.location.X()) + int(math.Round(dx))
	newYInt := int(c.location.Y()) + int(math.Round(dy))

	// Clamp to valid bounds
	if newXInt < int(kernel.MinX) {
		newXInt = int(kernel.MinX)
	} else if newXInt > int(kernel.MaxX) {
		newXInt = int(kernel.MaxX)
	}

	if newYInt < int(kernel.MinY) {
		newYInt = int(kernel.MinY)
	} else if newYInt > int(kernel.MaxY) {
		newYInt = int(kernel.MaxY)
	}

	newX := uint8(newXInt)
	newY := uint8(newYInt)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation
	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsInvalidError("orderID")
	}
	for _, place := range c.storagePlaces {
		if *(place.OrderID()) == orderID {
			return place, nil
		}
	}
	return nil, ErrNoOrderFound
}
