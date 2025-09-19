package orderrepo

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
)

func DomainToDTO(aggregate *order.Order) OrderDTO {
	var orderDTO OrderDTO
	orderDTO.ID = aggregate.ID()
	orderDTO.CourierID = aggregate.CourierID()
	orderDTO.Location = LocationDTO{
		X: int(aggregate.Location().X()),
		Y: int(aggregate.Location().Y()),
	}
	orderDTO.Volume = int(aggregate.Volume())
	orderDTO.Status = aggregate.Status()
	return orderDTO
}

func DtoToDomain(dto OrderDTO) *order.Order {
	var aggregate *order.Order
	location, _ := kernel.NewLocation(uint8(dto.Location.X), uint8(dto.Location.Y))
	aggregate = order.RestoreOrder(dto.ID, dto.CourierID, location, kernel.Volume(dto.Volume), dto.Status)
	return aggregate
}
