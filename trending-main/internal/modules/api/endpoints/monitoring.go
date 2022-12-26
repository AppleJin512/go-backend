package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"moonbite/trending/internal/blockchain"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"github.com/memclutter/gocore/pkg/coreslices"
)

func Monitoring(c echo.Context) error {
	ctx := c.Request().Context()
	instances := strings.Split(config.Config.Instances, ":")
	metrics := make([]schemas.MonitoringMetric, 0)

	// Solana nodes status
	i := 1
	for name, nodeSetting := range blockchain.NodeSettings {
		if !coreslices.StringIn(name, instances) {
			continue
		}
		ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
		rpcClient := rpc.New(nodeSetting.NodeUrlRpc)
		if _, err := rpcClient.GetHealth(ctx); err != nil {
			metrics = append(metrics, schemas.MonitoringMetric{
				Name:    fmt.Sprintf("Solana node[%s]", name),
				Status:  schemas.MetricStatusError,
				Message: err.Error(),
			})
		} else {
			metrics = append(metrics, schemas.MonitoringMetric{
				Name:   fmt.Sprintf("Solana node[%s]", name),
				Status: schemas.MetricStatusOk,
			})
		}
		cancel()
		i++
	}

	// Check cookie
	// Read cf pass cookie
	cfPass := schemas.AuthCookie{}
	cfPassCookie, err := models.Cache.Get(context.Background(), "cf_cookie").Bytes()
	if err != nil {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    fmt.Sprintf("Cloudflare cookie"),
			Status:  schemas.MetricStatusError,
			Message: fmt.Sprintf("redis.Get: %v", err),
		})
	} else if err == redis.Nil || len(cfPassCookie) == 0 {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    fmt.Sprintf("Cloudflare cookie"),
			Status:  schemas.MetricStatusError,
			Message: "not found cookie in redis, check chrome and cookies services",
		})
	} else if err := json.Unmarshal(cfPassCookie, &cfPass); err != nil {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    fmt.Sprintf("Cloudflare cookie"),
			Status:  schemas.MetricStatusError,
			Message: fmt.Sprintf("json.Unmarshall: %v", err),
		})
	} else {
		if ok, err := CheckCookie(cfPass); !ok {
			metrics = append(metrics, schemas.MonitoringMetric{
				Name:    fmt.Sprintf("Cloudflare cookie"),
				Status:  schemas.MetricStatusError,
				Message: fmt.Sprintf("check cookie error: %v", err),
			})
		} else {
			metrics = append(metrics, schemas.MonitoringMetric{
				Name:   "Cloudflare cookie",
				Status: schemas.MetricStatusOk,
				//Message: ,
			})
		}
	}

	// Collections in db
	collectionsCount, err := models.DB.NewSelect().Model((*models.Collection)(nil)).Count(ctx)
	if err != nil {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Collections in db",
			Status:  schemas.MetricStatusError,
			Message: fmt.Sprintf("db error: %v", err),
		})
	} else if collectionsCount <= 100 {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Collections in db",
			Status:  schemas.MetricStatusWarn,
			Message: fmt.Sprintf("%d, less than 100 collection", collectionsCount),
		})
	} else {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Collections in db",
			Status:  schemas.MetricStatusOk,
			Message: fmt.Sprintf("%d", collectionsCount),
		})
	}

	// Fresh activities (stats health)
	activitiesCount, err := models.DB.NewSelect().Model((*models.Activity)(nil)).Where("block_time >= ?", time.Now().UTC().Add(-1*time.Hour)).Count(ctx)
	if err != nil {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Stats in db per hour",
			Status:  schemas.MetricStatusError,
			Message: fmt.Sprintf("db error: %v", err),
		})
	} else if activitiesCount == 0 {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Stats in db per hour",
			Status:  schemas.MetricStatusWarn,
			Message: fmt.Sprintf("0, no stats per hour"),
		})
	} else {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Stats in db per hour",
			Status:  schemas.MetricStatusOk,
			Message: fmt.Sprintf("%d", activitiesCount),
		})
	}

	// Proxies
	proxies := pkg.GetProxies()
	total := len(proxies)
	live := 0
	dead := make([]string, 0)
	for _, proxyUrl := range proxies {
		if err := pkg.CheckProxy(proxyUrl); err == nil {
			live += 1
		} else {
			dead = append(dead, proxyUrl.Host)
		}
	}
	if len(dead) == 0 {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Proxies",
			Status:  schemas.MetricStatusOk,
			Message: fmt.Sprintf("%d", total),
		})
	} else if len(dead) < total/2 {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Proxies",
			Status:  schemas.MetricStatusWarn,
			Message: fmt.Sprintf("%d/%d, dead: %s", live, total, strings.Join(dead, ", ")),
		})
	} else {
		metrics = append(metrics, schemas.MonitoringMetric{
			Name:    "Proxies",
			Status:  schemas.MetricStatusError,
			Message: fmt.Sprintf("%d/%d, dead: %s", live, total, strings.Join(dead, ", ")),
		})
	}

	return c.JSON(http.StatusOK, metrics)
}
