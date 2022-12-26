package cli

import (
	"moonbite/trending/internal/cli/hooks"
	"moonbite/trending/internal/modules/actualize"
	"moonbite/trending/internal/modules/api"
	"moonbite/trending/internal/modules/collections"
	"moonbite/trending/internal/modules/cookies"
	"moonbite/trending/internal/modules/merarity"
	"moonbite/trending/internal/modules/stats"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:    "trending",
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "dsnDb", Value: "postgresql://postgres:root@localhost:5432/trending?sslmode=disable", EnvVars: []string{"DSN_DB"}},
			&cli.StringFlag{Name: "dsnCache", Value: "redis://:NBqthmT8RfuQt5Uh@localhost:6379/0", EnvVars: []string{"DSN_CACHE"}},
			&cli.BoolFlag{Name: "debug", Value: false, EnvVars: []string{"DEBUG"}},
			&cli.StringFlag{Name: "logLevel", Value: logrus.InfoLevel.String(), EnvVars: []string{"LOG_LEVEL"}},
			// First is operational!
			&cli.StringFlag{Name: "instances", Value: "GALAXITY_SUB", EnvVars: []string{"INSTANCES"}},
			&cli.StringFlag{Name: "centrifugoAddr", Value: "centrifugo", EnvVars: []string{"CENTRIFUGO_ADDR"}},
			&cli.StringFlag{Name: "centrifugoApiKey", Value: "", EnvVars: []string{"CENTRIFUGO_API_KEY"}},
		},
		Before: hooks.Before,
		After:  hooks.After,
		Commands: cli.Commands{
			{
				Name:   "api",
				Action: api.Command,
			},
			{
				Name:   "collections",
				Action: collections.Command,
			},
			{
				Name:   "stats",
				Action: stats.Command,
			},
			{
				Name:   "cookies",
				Action: cookies.Command,
			},
			{
				Name:   "actualize",
				Action: actualize.Command,
			},
			{
				Name:   "merarity",
				Action: merarity.Command,
			},
		},
	}
}
