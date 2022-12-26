package endpoints

import (
	"database/sql"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func CollectionsList(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(schemas.CollectionsListRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	collections := make([]models.Collection, 0)
	query := models.DB.NewSelect().Model(&collections)
	if req.Symbol != "" {
		query = query.Where("symbol = ?", req.Symbol)
	} else if req.MetaSymbol != "" {
		query = query.Where("meta_symbol = ?", req.MetaSymbol)
	} else if req.UpdateAuthority != "" {
		query = query.Where("update_authority = ?", req.UpdateAuthority)
	} else if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+strings.ToLower(req.Search)+"%")
	}
	if err := query.Limit(20).Scan(ctx); err == sql.ErrNoRows {
		return c.JSON(http.StatusOK, collections)
	} else if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, collections)
}
