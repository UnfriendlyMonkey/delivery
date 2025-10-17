package ports

import (
	"context"
	"delivery/internal/core/domain/kernel"
)

type GeoClient interface{
	GetGeoLocation(ctx context.Context, street string) (kernel.Location, error)
}
