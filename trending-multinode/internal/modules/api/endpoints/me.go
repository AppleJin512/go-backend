package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	. "github.com/memclutter/gorequests"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"
	"strings"
)

// Proxy method for magic
//
// https://stats-mainnet.magiceden.io/collection_stats/getCollectionTimeSeries/%7Bsymbol%7D
// https://api-mainnet.magiceden.dev/rpc/getNFTsByOwner/$%7Bwallet%7D
// https://api-mainnet.magiceden.dev/v2/collections/$%7Bsymbol%7D/stats
// https://api-mainnet.magiceden.dev/v2/tokens/$%7Bhash%7D
// https://api-mainnet.magiceden.io/volumes?edge_cache=true
// https://api-mainnet.magiceden.io/rpc/getCollectionHolderStats/$%7Bsymbol%7D?edge_cache=true
// https://api-mainnet.magiceden.dev/v2/wallets/$%7Bwallet%7D/activities?offset=0&limit=15
// https://api-mainnet.magiceden.dev/v2/tokens/$%7Bhash%7D

func ME(c echo.Context) error {
	ctx := c.Request().Context()
	path := strings.TrimPrefix(c.Request().RequestURI, "/me/")
	if strings.Contains(path, "api-mainnet.magiceden.dev") {
		logrus.WithFields(logrus.Fields{"path": path}).Info("proxy path dev")
		out := make([]byte, 0)
		if err := Get(
			WithExtensions(ProxiesExtension{Proxies: config.Proxies}),
			WithUrl("https://"+path),
			WithOut(&out, OutTypeBytes),
		); err != nil {
			return fmt.Errorf("error do get: %v", err)
		}
		return c.Blob(http.StatusOK, "application/json", out)
	} else {
		logrus.WithFields(logrus.Fields{"path": path}).Info("proxy path")
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

		url := fmt.Sprintf("https://" + path)

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
}
