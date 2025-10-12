// Package courierrepo
package courierrepo

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
)

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO
	courierDTO.ID = aggregate.ID()
	courierDTO.Name = aggregate.Name()
	courierDTO.Speed = aggregate.Speed()
	courierDTO.Location = LocationDTO{
		X: int(aggregate.Location().X()),
		Y: int(aggregate.Location().Y()),
	}
	storagePlaces := make([]*StoragePlaceDTO, 0, len(aggregate.StoragePlaces()))
	for _, sp := range aggregate.StoragePlaces() {
		spToDTO := SPDomainToDTO(sp)
		courierID := courierDTO.ID
		spToDTO.CourierID = &courierID
		storagePlaces = append(storagePlaces, &spToDTO)
	}
	courierDTO.StoragePlaces = storagePlaces
	return courierDTO
}

func SPDomainToDTO(entity *courier.StoragePlace) StoragePlaceDTO {
	var sPDTO StoragePlaceDTO
	sPDTO.ID = entity.ID()
	sPDTO.Name = entity.Name()
	sPDTO.TotalVolume = int(entity.TotalVolume())
	sPDTO.OrderID = entity.OrderID()
	return sPDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	var aggregate *courier.Courier
	location, _ := kernel.NewLocation(uint8(dto.Location.X), uint8(dto.Location.Y))
	storagePlaces := make([]*courier.StoragePlace, 0, len(dto.StoragePlaces))
	for _, sp := range dto.StoragePlaces {
		spToDomain := courier.RestoreStoragePlace(sp.Name, kernel.Volume(sp.TotalVolume), sp.ID, sp.OrderID)
		storagePlaces = append(storagePlaces, spToDomain)
	}
	aggregate = courier.RestoreCourier(dto.Name, dto.Speed, location, dto.ID, storagePlaces)
	return aggregate
}

func SPDtoToDomain(dto StoragePlaceDTO) *courier.StoragePlace {
	entity := courier.RestoreStoragePlace(dto.Name, kernel.Volume(dto.TotalVolume), dto.ID, dto.OrderID)
	return entity
}
