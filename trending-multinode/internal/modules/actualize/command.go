package actualize

import (
	"context"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/pkg"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Command(c *cli.Context) error {
	for {
		if err := Update(c.Context); err != nil {
			logrus.Errorf("error update collections: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}
		time.Sleep(10 * time.Minute)
	}
}

func Update(ctx context.Context) error {
	collections := make([]models.Collection, 0)
	if err := models.DB.NewSelect().Model(&collections).Where("update_authority IS NULL OR update_authority = ''").Scan(ctx); err != nil {
		return err
	}

	if len(collections) == 0 {
		logrus.Info("no blank")
		return nil
	}

	for i, collection := range collections {
		logrus.Infof("actualize %s", collection.Symbol)
		if err := pkg.GetCollectionData(ctx, config.RpcClient, collection.Symbol, &collections[i]); err != nil {
			logrus.Errorf("error get collection '%s' data: %v", collection.Name, err)
			continue
		}
	}

	if _, err := models.DB.NewUpdate().Model(&collections).Exec(ctx); err != nil {
		return err
	}

	return nil
}
