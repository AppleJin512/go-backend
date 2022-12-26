package endpoints

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"net/http"
)

func ListingsList(c echo.Context) (err error) {
	res := schemas.ListingsListResponse{}
	ctx := c.Request().Context()
	req := new(schemas.ListingsListRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	res.Items = make([]models.Listing, 0)
	query := models.DB.NewSelect().Model(&res.Items)
	if req.Symbol != "" {
		query = query.Where("symbol = ?", req.Symbol)
	}

	if req.Limit > 1000 || req.Limit <= 0 {
		req.Limit = 1000
	}
	query = query.Limit(req.Limit).Offset(req.Offset)

	if len(req.Order) > 0 {
		if req.Order[0] == '-' {
			req.Order = req.Order[1:] + " DESC"
		}
		query = query.OrderExpr(req.Order)
	} else {
		query = query.OrderExpr("block_time DESC")
	}

	if err = query.Scan(ctx); err == sql.ErrNoRows {
		return c.JSON(http.StatusOK, res)
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
