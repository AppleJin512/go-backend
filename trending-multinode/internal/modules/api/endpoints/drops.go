package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

func Drops(c echo.Context) error {

	ctx := c.Request().Context()
	reqData := new(schemas.DropsRequest)
	if err := c.Bind(reqData); err != nil {
		return err
	}

	if reqData.Limit == 0 {
		reqData.Limit = 250
	}

	cfPass := schemas.AuthCookie{}
	cfPassCookie, err := models.Cache.Get(ctx, "cf_cookie").Bytes()
	if err != nil {
		return err
	} else if err == redis.Nil || len(cfPassCookie) == 0 {
		return fmt.Errorf("nil in cf")
	} else if err := json.Unmarshal(cfPassCookie, &cfPass); err != nil {
		return err
	}

	client, err := pkg.ByPassCloudFlareClient()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api-mainnet.magiceden.io/drops?limit=%d&offset=%d", reqData.Limit, reqData.Offset)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error create request: %v", err)
	}
	req.AddCookie(&http.Cookie{
		Name:     "__cf_bm",
		Value:    cfPass.Cookie,
		Path:     "/",
		Domain:   ".magiceden.io",
		Secure:   true,
		HttpOnly: true,
	})
	req.Header.Set("sec-ch-ua", cfPass.SecShUa)
	req.Header.Set("user-agent", cfPass.UserAgent)
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error get do: %v", err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(data) == 0 || data[0] != '[' {
		return c.JSON(http.StatusInternalServerError, schemas.Error{
			Message: "Error server",
		})
	}

	return c.Blob(http.StatusOK, "application/json", data)
}
