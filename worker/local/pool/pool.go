package pool

import (
	"context"
)

type IWorker interface {
	ID() string
	Execute(ctx context.Context, params []byte) (any, error)
}
