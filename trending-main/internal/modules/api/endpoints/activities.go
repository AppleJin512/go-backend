package endpoints

import (
	"database/sql"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/memclutter/gocore/pkg/coreslices"
	"github.com/uptrace/bun"
)

func ActivitiesList(c echo.Context) (err error) {
	res := schemas.ActivitiesListResponse{}
	ctx := c.Request().Context()
	req := new(schemas.ActivitiesListRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	res.Items = make([]models.Activity, 0)
	query := models.DB.NewSelect().Model(&res.Items).Relation("Item")
	if req.Symbol != "" {
		query = query.Where("symbol = ?", req.Symbol)
	}

	activityTypes := strings.Split(req.ActivityType, ",")
	activityTypes = coreslices.StringApply(activityTypes, func(i int, s string) string { return strings.TrimSpace(s) })
	activityTypes = coreslices.StringFilter(activityTypes, func(i int, s string) bool { return len(s) > 0 })
	if len(activityTypes) > 0 {
		query = query.Where("activity_type IN (?)", bun.In(activityTypes))
	}

	if req.Limit > 1000 || req.Limit <= 0 {
		req.Limit = 1000
	}
	query = query.Limit(req.Limit).Offset(req.Offset).OrderExpr("block_time DESC")

	if err = query.Scan(ctx); err == sql.ErrNoRows {
		return c.JSON(http.StatusOK, res)
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
