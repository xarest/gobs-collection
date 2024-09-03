package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/lib/config"
	"github.com/xarest/gobs-template/lib/logger"
	gCommon "github.com/xarest/gobs/common"
)

type DBConfig struct {
	Uri      string `env:"DB_URI" mapstructure:"DB_URI" envDefault:""`
	Type     string `env:"DB_TYPE" mapstructure:"DB_TYPE" envDefault:"postgres"`
	Host     string `env:"DB_HOST" mapstructure:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" mapstructure:"DB_PORT" envDefault:"5432"`
	DbName   string `env:"DB_NAME" mapstructure:"DB_NAME" envDefault:"postgres"`
	UserName string `env:"DB_USER" mapstructure:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" mapstructure:"DB_PASSWORD" envDefault:"postgres"`
	SslMode  string `env:"DB_SSL_MODE" mapstructure:"DB_SSL_MODE" envDefault:"disable"`
}

type DB struct {
	*bun.DB

	log logger.ILogger
}

func (d *DB) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			config.NewIConfig(),
		},
		AsyncMode: map[gCommon.ServiceStatus]bool{
			gCommon.StatusSetup: true,
			gCommon.StatusStart: true,
			gCommon.StatusStop:  true,
		},
	}, nil
}

func (d *DB) Setup(ctx context.Context, deps gobs.Dependencies) error {
	var (
		sCfg     config.IConfiguration
		dbConfig DBConfig
	)
	if err := deps.Assign(&d.log, &sCfg); err != nil {
		return err
	}

	if err := sCfg.Parse(&dbConfig); err != nil {
		return err
	}
	if dbConfig.Uri == "" {
		// dbConfig.Uri = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		// 	dbConfig.Host, dbConfig.Port, dbConfig.UserName, dbConfig.Password, dbConfig.DbName, dbConfig.SslMode,
		// )
		dbConfig.Uri = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName, dbConfig.SslMode)
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbConfig.Uri)))
	d.DB = bun.NewDB(sqldb, pgdialect.New())

	return nil
}

func (d *DB) Start(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}

func (d *DB) Stop(ctx context.Context) error {
	return d.DB.Close()
}

var _ gobs.IServiceInit = (*DB)(nil)
var _ gobs.IServiceSetup = (*DB)(nil)
var _ gobs.IServiceStart = (*DB)(nil)
var _ gobs.IServiceStop = (*DB)(nil)
