package services

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)


func Test_OrderDispatcherServiceErrorNoOrder(t *testing.T) {
	// Arrange
	var o *order.Order
	couriers := []*courier.Courier{
		courier.CreateCourierOK(),
	}

	// Act
	orderDispatcherService := NewOrderDispatcherService()
	c, err := orderDispatcherService.Dispatch(o, couriers)

	// Assert
	assert.Error(t, err, "should be error dispatching with no order")
	assert.Nil(t, c, "no courier should be returned when dispatching with no order")
}

func Test_OrderDispatcherServiceErrorNoCouriers(t *testing.T) {
	// Arrange
	var o = order.CreateOrderOK()
	couriers := []*courier.Courier{}

	// Act
	orderDispatcherService := NewOrderDispatcherService()
	c, err := orderDispatcherService.Dispatch(o, couriers)

	// Assert
	assert.Error(t, err, "should be error dispatching without couriers")
	assert.Nil(t, c, "no courier should be returned when dispatching without couriers")
}

func Test_OrderDispatcherServiceOkOneCourier(t *testing.T) {
	// Arrange
	var o = order.CreateOrderOK()
	couriers := []*courier.Courier{
		courier.CreateCourierOK(),
	}

	// Act
	orderDispatcherService := NewOrderDispatcherService()
	c, err := orderDispatcherService.Dispatch(o, couriers)

	// Assert
	assert.NoError(t, err, "should be no error dispatching wiht correct params")
	assert.NotEmpty(t, c, "courier should be returned when dispatching with correct params")
	assert.True(t, c.Equal(couriers[0]), "winner should be equal to single participant")
	assert.Equal(t, order.StatusAssigned, o.Status(), "order status should be Assigned after dispatching")
}

func Test_OrderDispatcherServiceOkManyCouriers(t *testing.T) {
	// Arrange
	var o = order.CreateOrderOK()
	couriers := []*courier.Courier{
		courier.CreateCourierOK(),
		courier.CreateCourierOK(),
		courier.CreateCourierOK(),
		courier.CreateCourierOK(),
	}

	// Act
	orderDispatcherService := NewOrderDispatcherService()
	c, err := orderDispatcherService.Dispatch(o, couriers)

	// Assert
	assert.NoError(t, err, "should be no error dispatching wiht correct params")
	assert.NotEmpty(t, c, "courier should be returned when dispatching with correct params")
}

func Test_OrderDispatcherServiceErrorCourierBuisy(t *testing.T) {
	// Arrange
	var o = order.CreateOrderOK()
	cc := courier.CreateCourierOK()
	couriers := []*courier.Courier{
		cc,
	}

	// Act
	_ = cc.TakeOrder(order.CreateOrderOK())
	orderDispatcherService := NewOrderDispatcherService()
	c, err := orderDispatcherService.Dispatch(o, couriers)

	// Assert
	assert.Error(t, err, "should be error dispatching with buisy couriers")
	assert.Equal(t, ErrCourierNotFound, err, fmt.Sprintf("expected %v, got %v", ErrCourierNotFound, err))
	assert.Nil(t, c, "no courier should be returned when dispatching with buisy couriers")
}

func Test_OrderDispatcherServiceBestTime(t *testing.T) {
	o, _ := order.NewOrder(uuid.New(), kernel.MinLocation(), kernel.Volume(kernel.MinVolume))
	c1, _ := courier.NewCourier("one", 1, kernel.MaxLocation())
	c2, _ := courier.NewCourier("two", 2, kernel.MaxLocation())
	c3, _ := courier.NewCourier("three", 4, kernel.MaxLocation())
	couriers := []*courier.Courier{
		c1, c2, c3,
	}

	// Act
	orderDispatcherService := NewOrderDispatcherService()
	c, err := orderDispatcherService.Dispatch(o, couriers)
	
	// Assert
	assert.NoError(t, err, "should be no error dispatching wiht correct params")
	assert.NotEmpty(t, c, "courier should be returned when dispatching with correct params")
	assert.Equal(t, c3, c, "fastest courier should be returned")
}
