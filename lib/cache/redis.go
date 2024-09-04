package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/lib/config"
	"github.com/xarest/gobs-template/lib/logger"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type Redis struct {
	rClient *redis.Client
	log     logger.ILogger
}

func (o *Redis) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: []gobs.IService{
			config.NewIConfig(),
			logger.NewILogger(),
		},
	}, nil
}

func (o *Redis) Setup(c context.Context, deps ...gobs.IService) error {
	var (
		rdbCfg RedisConfig
		config config.IConfiguration
	)
	if err := gobs.Dependencies(deps).Assign(&config, &o.log); err != nil {
		return err
	}
	if err := config.Parse(&rdbCfg); err != nil {
		return err
	}
	o.rClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdbCfg.Host, rdbCfg.Port), // Redis server address
		Password: rdbCfg.Password,                                // No password by default
		DB:       rdbCfg.DB,                                      // Default DB
	})
	return nil
}

func (o *Redis) Start(ctx context.Context) error {
	_, err := o.rClient.Ping(ctx).Result()
	return err
}

func (o *Redis) Stop(ctx context.Context) error {
	return o.rClient.Close()
}

func (r *Redis) Set(c context.Context, key string, value any, ttls time.Duration) error {
	rawData, err := json.Marshal(value)
	if nil != err {
		return err
	}

	return r.rClient.Set(c, key, rawData, ttls).Err()
}

func (r *Redis) Get(c context.Context, key string, model any) (any, error) {
	rawData, err := r.rClient.Get(c, key).Result()
	if nil != err {
		return nil, err
	}
	if len(rawData) == 0 {
		return nil, errors.New("data empty")
	}
	return model, json.Unmarshal([]byte(rawData), model)
}

func (r *Redis) Wrap(
	c context.Context,
	key string,
	model any,
	callback func() (any, error),
	ttls time.Duration,
) (res any, err error) {
	res, err = r.Get(c, key, model)
	if err != nil {
		res, err := callback()
		if nil != err {
			return nil, err
		}
		return res, r.Set(c, key, res, ttls)
	}
	return res, err
}

var _ gobs.IServiceInit = (*Redis)(nil)
var _ gobs.IServiceSetup = (*Redis)(nil)
var _ gobs.IServiceStart = (*Redis)(nil)
var _ gobs.IServiceStop = (*Redis)(nil)
var _ ICache = (*Redis)(nil)
