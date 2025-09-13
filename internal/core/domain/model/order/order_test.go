package order_test

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	VolumeOK = 5
)

func CreateOrderOK() *order.Order {
	orderID := uuid.New()
	location, _ := kernel.RandomLocation()
	volume, _ := kernel.NewVolume(VolumeOK)
	o, _ := order.NewOrder(orderID, location, *volume)
	return o
}

func Test_NewOrderOkWithValidParams(t *testing.T) {
	// Arrange
	orderID := uuid.New()
	locations, err := kernel.RandomLocation()
	assert.NoError(t, err, "should be no error creating random location")
	volume, err := kernel.NewVolume(VolumeOK)
	assert.NoError(t, err, "should be no error creating new volume")

	// Act
	o, err := order.NewOrder(orderID, locations, *volume)

	// Assert
	assert.NoError(t, err, "should be no error creating Order with valid params")
	assert.NotEmpty(t, o, "order should not be nil")
	assert.Equal(t, orderID, o.ID(), "order ID should match input param")
	assert.Nil(t, o.CourierID(), "new order should not have courierID")
	assert.Equal(t, order.StatusCreated, o.Status(), "new order status shuld be 'Created'")
}

func Test_NewOrderErrorsWithWrongParams(t *testing.T) {
	// Arrange
	okLocation, _ := kernel.RandomLocation()
	okVolume, _ := kernel.NewVolume(VolumeOK)
	tests := map[string]struct {
		orderID uuid.UUID
		location kernel.Location
		volume kernel.Volume
		expected error
	}{
		"wrong_id": {
			orderID: uuid.Nil,
			location: okLocation,
			volume: *okVolume,
			expected: errs.NewValueIsInvalidError("orderID"),
		},
		"wrong_location": {
			orderID: uuid.New(),
			location: kernel.Location{},
			volume: *okVolume,
			expected: errs.NewValueIsInvalidError("location"),
		},
		"wrong_volume": {
			orderID: uuid.New(),
			location: okLocation,
			volume: kernel.Volume(0),
			expected: errs.NewValueIsInvalidError("volume"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := order.NewOrder(test.orderID, test.location, test.volume)
			assert.Equal(t, test.expected, err, fmt.Sprintf("expected %v, got %v", test.expected, err))
		})
	}
}

func Test_OrderCanAssignCourierOK(t *testing.T) {
	o := CreateOrderOK()
	courierID := uuid.New()
	err := o.Assign(&courierID)
	assert.NoError(t, err, "should be no error assigning client to order")
	assert.Equal(t, o.CourierID(), &courierID, "order's courier ID should match courierID")
	assert.Equal(t, order.StatusAssigned, o.Status(), "order status should be 'Assigned' after Assignment")
}

func Test_OrderAssignErrorWrongID(t *testing.T) {
	targetErr := errs.NewValueIsInvalidError("courierID")
	o := CreateOrderOK()
	courierID := uuid.Nil
	err := o.Assign(&courierID)
	assert.Error(t, err, "should be arror assigning empty courier ID to order")
	assert.Equal(t, err, targetErr, fmt.Sprintf("expected %v, got %v", targetErr, err))
}

func Test_OrderCompleteOK(t *testing.T) {
	o := CreateOrderOK()
	courierID := uuid.New()
	_ = o.Assign(&courierID)
	err := o.Complete()
	assert.NoError(t, err, "should be no error compliting assigned order")
	assert.Equal(t, order.StatusCompleted, o.Status(), "order status should be 'Completed' after completion")
}

func Test_OrderCompleteErrorNotAssigned(t *testing.T) {
	o := CreateOrderOK()
	err := o.Complete()
	assert.Error(t, err, "should be error completing not assigned order")
	assert.Equal(t, order.ErrOrderStatusIsWrongForAction, err, fmt.Sprintf(
		"expected %v, got %v", order.ErrOrderStatusIsWrongForAction, err))
}
