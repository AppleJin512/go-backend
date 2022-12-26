package collections

import (
	"context"
	"fmt"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/pkg"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/memclutter/gorequests"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Command(c *cli.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func(c *cli.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := InitCommand(c); err != nil {
			logrus.Errorf("error init command: %v", err)
		}
	}(c, wg)

	go func(c *cli.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := MainCommand(c); err != nil {
			logrus.Errorf("error main command: %v", err)
		}
	}(c, wg)

	wg.Wait()
	return nil
}

func MainCommand(c *cli.Context) error {
	for {
		collections, err := GetPage(0, 10)
		if err != nil {
			logrus.Errorf("error get new collections: %v", err)
			time.Sleep(2 * time.Minute)
			continue
		}
		if len(collections) > 0 {
			if err := Save(c.Context, collections); err != nil {
				logrus.Errorf("error save new collections: %v", err)
				time.Sleep(2 * time.Minute)
				continue
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func InitCommand(c *cli.Context) error {
	offset := 0
	limit := 500

	for {
		collections, err := GetPage(offset, limit)
		if err != nil {
			logrus.Errorf("error get page %d: %v", offset, err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Finish by empty page
		if len(collections) == 0 {
			break
		}

		if err := Save(c.Context, collections); err != nil {
			logrus.Errorf("error save page %d: %v", offset, err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Finish by last page
		if len(collections) < limit {
			break
		}

		// Next
		offset += limit
	}

	logrus.Infof("init complete")
	return nil
}

func Save(ctx context.Context, raw []MagicCollection) error {
	collections := make([]models.Collection, 0)
	for _, r := range raw {
		r.Symbol = strings.TrimSpace(r.Symbol)
		if len(r.Symbol) == 0 {
			logrus.Warnf("empty symbol")
			continue
		}

		collection := models.Collection{
			Symbol: r.Symbol,
			Name:   r.Name,
			Image:  r.Image,
		}

		logrus.Infof("actualize %s", collection.Symbol)
		if err := pkg.GetCollectionData(context.Background(), config.RpcClient, collection.Symbol, &collection); err != nil {
			logrus.Errorf("error get collection '%s' data: %v", collection.Name, err)
			continue
		}
		collections = append(collections, collection)
		time.Sleep(200 * time.Millisecond)
	}

	if _, err := models.DB.NewInsert().Model(&collections).On("conflict(symbol) DO UPDATE").Set("meta_symbol = EXCLUDED.meta_symbol").Set("update_authority = EXCLUDED.update_authority").Exec(ctx); err != nil {
		return fmt.Errorf("error insert page: %v", err)
	}

	return nil
}

func GetPage(offset, limit int) ([]MagicCollection, error) {
	logrus.Infof("get page %d:%d", offset, limit)
	collections := make([]MagicCollection, 0)
	if err := gorequests.Get(
		gorequests.WithExtensions(
			gorequests.ProxiesExtension{Proxies: config.Proxies},
			gorequests.RetryExtension{RetryMax: 5, RetryWaitMin: 500 * time.Millisecond, RetryWaitMax: 800 * time.Millisecond},
		),
		gorequests.WithUrl("https://api-mainnet.magiceden.dev/v2/collections?offset=%d&limit=%d&order=date", offset, limit),
		gorequests.WithOkStatusCodes(http.StatusOK),
		gorequests.WithOut(&collections, gorequests.OutTypeJson),
	); err != nil {
		return collections, fmt.Errorf("gorequests get error: %v", err)
	}
	return collections, nil
}
