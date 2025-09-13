package courier_test

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"fmt"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	NameOK = "SomeName"
	NameEmpty = ""
	SpeedOK = 2
	SpeedLow = 0
	SpeedHigh = 9
)

func createCourierOK() *courier.Courier {
	location, _ := kernel.RandomLocation()
	c, _ := courier.NewCourier(NameOK, SpeedOK, location)
	return c
}

func CreateOrderOK() *order.Order {
	// TODO: move to order package
	orderID := uuid.New()
	location, _ := kernel.RandomLocation()
	volume, _ := kernel.NewVolume(VolumeOK)
	o, _ := order.NewOrder(orderID, location, *volume)
	return o
}

func Test_NewCourierOkWithValidParams(t *testing.T) {
	// Arrange
	location, err := kernel.RandomLocation()
	assert.NoError(t, err, "should be no error creating random location")

	// Act
	c, err := courier.NewCourier(NameOK, SpeedOK, location)

	// Assert
	assert.NoError(t, err, "should be no error creating Courier with valid params")
	assert.NotEmpty(t, c, "new courier should not be nil")
	assert.Equal(t, NameOK, c.Name(), "name should match input param")
	assert.Equal(t, SpeedOK, c.Speed(), "speed shuld match input param")
	assert.Equal(t, 1, len(c.StoragePlaces()), "new courier must have one storage place")
	assert.NotEmpty(t, c.ID(), "new courier Id should not be empty")
}

func Test_NewCourierErrorsWithWrongParams(t *testing.T) {
	// Arrange
	okLocation, _ := kernel.RandomLocation()
	tests := map[string]struct {
		name string
		speed int
		location kernel.Location
		expected error
	}{
		"wrong_name": {
			name: NameEmpty,
			speed: SpeedOK,
			location: okLocation,
			expected: errs.NewValueIsInvalidError("name"),
		},
		"wrong_location": {
			name: NameOK,
			speed: SpeedOK,
			location: kernel.Location{},
			expected: errs.NewValueIsInvalidError("location"),
		},
		"wrong_speed_too_high": {
			name: NameOK,
			speed: SpeedHigh,
			location: okLocation,
			expected: errs.NewValueIsOutOfRangeError("speed", SpeedHigh, courier.MinSpeed, courier.MaxSpeed),
		},
		"wrong_speed_too_low": {
			name: NameOK,
			speed: SpeedLow,
			location: okLocation,
			expected: errs.NewValueIsOutOfRangeError("speed", SpeedLow, courier.MinSpeed, courier.MaxSpeed),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := courier.NewCourier(test.name, test.speed, test.location)
			assert.Equal(t, test.expected, err, fmt.Sprintf("expected %v, got %v", test.expected, err))
		})
	}
}

func Test_CourierEqualOK(t *testing.T) {
	c := createCourierOK()
	assert.True(t, c.Equal(c), "courier should be equal to themselfs")
}

func Test_CourierEqualNot(t *testing.T) {
	c := createCourierOK()
	c2 := createCourierOK()
	assert.False(t, c.Equal(c2), "courier should not be equal other courier")
}

func Test_CourierCanAddStoragePlaceOK(t *testing.T) {
	c := createCourierOK()
	err := c.AddStoragePlace("pocket", courier.BagVolume)
	assert.NoError(t, err, "should be no error adding normal storage place")
	assert.Equal(t, 2, len(c.StoragePlaces()), "should be 2 storage places after adding")
}

func Test_CourierCanNotAddStoragePlaceWrongSize(t *testing.T) {
	c := createCourierOK()
	err := c.AddStoragePlace("empty_pocket", 0)
	assert.Error(t, err, "should be error adding too small storage place")
}

func Test_CourierCanNotAddStoragePlaceWrongName(t *testing.T) {
	c := createCourierOK()
	err := c.AddStoragePlace("", courier.BagVolume)
	assert.Error(t, err, "should be error adding storage place with no name")
}

func Test_CourierCalculateTimeOK(t *testing.T) {
	location := kernel.MinLocation()
	l2 := kernel.MaxLocation()
	c, _ := courier.NewCourier(NameOK, SpeedOK, location)
	dist, _ := location.Distance(l2)
	speed := c.Speed()
	expectedTime := math.Ceil(float64(dist) / float64(speed))
	time, err := c.CalculateTimeToLocation(l2)
	assert.NoError(t, err, "should calculate right between correct locations")
	assert.Equal(t, expectedTime, time, "time should match expected")
}

func Test_CourierTakeOrderOK(t *testing.T) {
	c := createCourierOK()
	o := CreateOrderOK()
	err := c.TakeOrder(o)
	assert.NoError(t, err, "courier should take normal order")
}

func Test_CourierCompleteOrderOK(t *testing.T) {
	c := createCourierOK()
	o := CreateOrderOK()
	err := c.TakeOrder(o)
	assert.NoError(t, err, "courier should take normal order")
	err = c.CompleteOrder(o)
	assert.NoError(t, err, "completing correct order should be OK")
}
