package endpoints

import (
	"database/sql"
	"fmt"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/memclutter/gocore/pkg/coreslices"
)

func Stats(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(schemas.StatsRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	periodDuration, err := time.ParseDuration(req.Period)
	if err != nil {
		return c.JSON(http.StatusBadRequest, schemas.Error{
			Message: fmt.Sprintf("Invalid period value %s, support values like 5m, 1h, 1d ({digit}{unit: m, h, d})", req.Period),
		})
	}

	blockTimeStart := time.Now().UTC().Add(-periodDuration)
	query := `		
		WITH stats AS (
			SELECT symbol,
				   count(*)   AS stats_trade_count,
				   round(avg(price), 2) AS stats_average_price,
				   round(sum(price), 2) AS stats_volume
			FROM activities
			WHERE activity_type = 'sale' 
              AND block_time BETWEEN ? AND NOW()
			GROUP BY symbol
		)
		SELECT c.name,
		       c.symbol,
		       c.image,
		       coalesce(s.stats_trade_count, 0),
		       coalesce(s.stats_average_price, 0), 
		       coalesce(s.stats_volume, 0)
		FROM collections AS c
		INNER JOIN stats AS s USING (symbol)
		ORDER BY s.stats_volume DESC
		LIMIT 20
	`

	// Replace default order if exists in request and translate string 'field,-field' to 'ORDER BY field, field DESC'
	req.Sort = strings.TrimSpace(req.Sort)
	if len(req.Sort) > 0 {
		sortSlice := strings.Split(req.Sort, ",")
		sortOrder := make([]string, 0)
		for i := range sortSlice {
			desc := ""
			field := strings.TrimSpace(sortSlice[i])
			if len(field) == 0 {
				continue
			} else if field[0] == '-' {
				field = field[1:]
				desc = " DESC"
			}
			if !coreslices.StringIn(field, []string{"stats_volume", "stats_trade_count", "stats_average_price", "name", "symbol"}) {
				continue
			}
			sortOrder = append(sortOrder, field+desc)
		}
		if len(sortOrder) > 0 {
			query = strings.Replace(query, "ORDER BY s.stats_volume DESC", "ORDER BY "+strings.Join(sortOrder, ", "), -1)
		}
	}

	rows, err := models.DB.QueryContext(ctx, query, blockTimeStart)
	if err != nil {
		return err
	}
	defer rows.Close()
	result := make([]schemas.StatsResponse, 0)
	for rows.Next() {
		item := schemas.StatsResponse{}
		if err := rows.Scan(
			&item.Name,
			&item.Symbol,
			&item.Image,
			&item.StatsTradeCount,
			&item.StatsAveragePrice,
			&item.StatsVolume,
		); err != nil {
			return err
		}
		result = append(result, item)
	}
	return c.JSON(http.StatusOK, result)
}

func StatsDistribution(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(schemas.StatsDistributionRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	periodDuration, err := time.ParseDuration(req.Period)
	if err != nil {
		return c.JSON(http.StatusBadRequest, schemas.Error{
			Message: fmt.Sprintf("Invalid period value %s, support values like 5m, 1h ({digit}{unit: m, h})", req.Period),
		})
	}

	//	SELECT to_timestamp(floor((extract('epoch' FROM block_time) / 600)) * 600)
	//				   AT TIME ZONE 'UTC' AS slot,
	//			   COUNT(*) FILTER (WHERE activity_type IN ('sale')) AS                  sale_cnt,
	//			   COUNT(*) FILTER (WHERE activity_type IN ('listing', 'update_price')) AS listing_cnt
	//		FROM activities
	//		WHERE symbol = ?
	//		GROUP BY 1
	//		ORDER BY 1 DESC
	res := make([]schemas.StatsDistributionResponse, 0)
	query := models.DB.NewSelect().Model(&res).
		ColumnExpr(`to_timestamp(floor((extract('epoch' FROM block_time) / ?)) * ?) AT TIME ZONE 'UTC' AS slot`, periodDuration.Seconds(), periodDuration.Seconds()).
		ColumnExpr(`COUNT(*) FILTER (WHERE activity_type IN ('sale')) AS count_sale`).
		ColumnExpr(`COUNT(*) FILTER (WHERE activity_type IN ('listing', 'update_price')) AS count_listing`).
		ModelTableExpr("activities").
		GroupExpr("1").OrderExpr("1 DESC")

	if len(req.Symbol) > 0 {
		query = query.Where("symbol = ?", req.Symbol)
	}

	if req.Limit == 0 || req.Limit > 12 {
		req.Limit = 12
	}
	query = query.Limit(req.Limit)

	if err := query.Scan(ctx); err == sql.ErrNoRows {
		return c.JSON(http.StatusOK, res)
	} else if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
