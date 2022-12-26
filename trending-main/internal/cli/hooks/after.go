package hooks

import (
	"moonbite/trending/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func After(c *cli.Context) error {
	defer func() {
		if models.DB == nil {
			return
		}
		if err := models.DB.Close(); err != nil {
			logrus.Errorf("db close error: %v", err)
		}
	}()
	return nil
}
