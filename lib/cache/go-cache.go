package cache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/xarest/gobs"
)

type GoCache struct {
	mc *cache.Cache
}

func (o *GoCache) Setup(c context.Context, deps gobs.Dependencies) error {
	o.mc = cache.New(5*time.Minute, 10*time.Minute)
	return nil
}

func (r *GoCache) Set(c context.Context, key string, value any, ttls time.Duration) error {
	r.mc.Set(key, value, ttls)
	return nil
}

func (r *GoCache) Get(c context.Context, key string, _ any) (any, error) {
	res, found := r.mc.Get(key)
	if !found {
		return nil, errors.New("key not found")
	}
	return res, nil
}

func (r *GoCache) Wrap(
	c context.Context,
	key string,
	_ any,
	callback func() (any, error),
	ttls time.Duration,
) (res any, err error) {
	res, err = r.Get(c, key, nil)
	if err != nil {
		res, err = callback()
		if nil != err {
			return res, err
		}
		return res, r.Set(c, key, res, ttls)
	}
	return res, nil
}

var _ gobs.IServiceSetup = (*Redis)(nil)
var _ ICache = (*GoCache)(nil)
