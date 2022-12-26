package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

func MESeries(c echo.Context) error {

	ctx := c.Request().Context()
	reqData := new(schemas.MESeriesRequest)
	if err := c.Bind(reqData); err != nil {
		return err
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

	url := fmt.Sprintf("https://stats-mainnet.magiceden.io/collection_stats/getCollectionTimeSeries/%s?edge_cache=true&resolution=%s&addLastDatum=true", reqData.Collection, reqData.Period)

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
	//req.Header.Set("sec-ch-ua", cfPass.SecShUa)
	//req.Header.Set("user-agent", cfPass.UserAgent)
	logrus.WithFields(logrus.Fields{
		"cookie": cfPass.Cookie,
		//"sec-sh-ua":  cfPass.SecShUa,
		//"user-agent": cfPass.UserAgent,
	}).Info("cloudflare bypass cookie")
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error get do: %v", err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/json", data)
}
