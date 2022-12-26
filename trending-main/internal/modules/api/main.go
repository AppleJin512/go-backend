package api

import (
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/modules/api/endpoints"
	"moonbite/trending/internal/modules/api/middleware"
	"moonbite/trending/internal/modules/api/plugins"

	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

func Command(c *cli.Context) error {
	e := echo.New()
	e.Debug = config.Config.Debug

	corsMiddleware := middleware.NewCors()

	e.Use(corsMiddleware)

	e.Validator = plugins.NewValidator()
	e.HTTPErrorHandler = plugins.NewErrorHandler(e)

	e.GET("/listings/", endpoints.ListingsList)
	e.GET("/stats/", endpoints.Stats)
	e.GET("/stats/distribution/", endpoints.StatsDistribution)
	e.GET("/activities/", endpoints.ActivitiesList)
	e.GET("/me-series/", endpoints.MESeries)
	e.GET("/drops/", endpoints.Drops)
	e.GET("/collections/", endpoints.CollectionsList)
	e.POST("/auth/cookie/", endpoints.AuthCookie)
	e.GET("/me-collection/", endpoints.MECollectionData)
	e.GET("/me-wallet/", endpoints.MEWalletData)
	e.GET("/me-stats/", endpoints.MEStats)
	e.GET("/me-buy-now/", endpoints.MEBuyNow)
	e.GET("/me-attributes/", endpoints.MeAttributes)
	e.GET("/monitoring/", endpoints.Monitoring)
	e.GET("/me/*", endpoints.ME)

	return e.Start(":8000")
}
