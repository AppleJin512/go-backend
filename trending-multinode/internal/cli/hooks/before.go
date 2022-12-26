package hooks

import (
	"database/sql"
	"fmt"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/models"

	"github.com/centrifugal/gocent/v3"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/urfave/cli/v2"
)

func Before(c *cli.Context) error {
	config.Config.Debug = c.Bool("debug")
	config.Config.Instances = c.String("instances")

	if err := config.Config.SetLogLevel(c.String("logLevel")); err != nil {
		return fmt.Errorf("cli hooks before error: %v", err)
	}
	logrus.SetLevel(config.Config.LogLevel)
	models.DB = bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.String("dsnDb")))), pgdialect.New())

	// For debug mode, enable sql queries log and debug log level
	if config.Config.Debug {
		models.DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		if config.Config.LogLevel < logrus.DebugLevel {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}
	var err error
	config.RpcClient, err = config.InitRpcClient(c.String("instances"))
	if err != nil {
		return err
	}

	optRd, err := redis.ParseURL(c.String("dsnCache"))
	if err != nil {
		return fmt.Errorf("error parse cache dsn: %v", err)
	}
	models.Cache = redis.NewClient(optRd)

	config.GoCent = gocent.New(gocent.Config{
		Addr: c.String("centrifugoAddr"),
		Key:  c.String("centrifugoApiKey"),
	})

	return nil
}
