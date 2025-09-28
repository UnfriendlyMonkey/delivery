package ports

import "context"

type UnitOfWorkFactory interface {
	New(ctx context.Context) (UnitOfWork, error)
}
