package cache

import (
	"context"
	"time"
)

type ICache interface {
	Set(c context.Context, key string, value any, ttls time.Duration) error
	Get(c context.Context, key string, model any) (any, error)
	Wrap(c context.Context, key string, model any, callback func() (any, error), ttls time.Duration) (any, error)
}

func NewICache() ICache {
	return &Redis{}
}
