package config

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/xarest/gobs"
)

type EnvConfig struct {
}

func (c *EnvConfig) Setup(ctx context.Context, _ gobs.Dependencies) error {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		// return err
	}
	return nil
}

func (c *EnvConfig) Parse(result interface{}) error {
	return env.Parse(result)
}

var _ IConfiguration = (*EnvConfig)(nil)
var _ gobs.IServiceSetup = (*EnvConfig)(nil)
