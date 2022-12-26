package endpoints

import (
	"fmt"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/modules/api/schemas"
	"net/http"

	"github.com/labstack/echo/v4"
	. "github.com/memclutter/gorequests"
)

func MECollectionData(c echo.Context) error {
	out := make([]byte, 0)
	req := schemas.MECollectionDataRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	} else if err := Get(
		WithExtensions(ProxiesExtension{Proxies: config.Proxies}),
		WithUrl("https://api-mainnet.magiceden.dev/collections/%s?edge_cache=true", req.Collection),
		WithOkStatusCodes(http.StatusOK),
		WithOut(&out, OutTypeBytes),
	); err != nil {
		return fmt.Errorf("error do get: %v", err)
	}
	return c.Blob(http.StatusOK, "application/json", out)
}
