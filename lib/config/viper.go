package config

import (
	"context"

	"github.com/spf13/viper"
	"github.com/xarest/gobs"
)

type Viper struct {
}

func (c *Viper) Setup(ctx context.Context, _ gobs.Dependencies) error {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	return viper.ReadInConfig()
}

func (c *Viper) Parse(result interface{}) error {
	return viper.Unmarshal(result)
}

var _ IConfiguration = (*Viper)(nil)
var _ gobs.IServiceSetup = (*Viper)(nil)
