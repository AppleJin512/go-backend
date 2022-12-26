package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func MeAttributes(c echo.Context) error {
	ctx := c.Request().Context()
	req := schemas.MeAttributesRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	} else if err := c.Validate(&req); err != nil {
		return err
	}
	attributes := make([]models.Attribute, 0)

	// Check cache data
	logrus.WithField("collection", req.Collection).Infof("select from db")
	query := models.DB.NewSelect().Model(&attributes).
		Where("symbol = ?", req.Collection).
		Where("date_updated > ?", time.Now().UTC().Add(-24*time.Hour))
	if err := query.Scan(ctx); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("select db error: %v", err)
	} else if len(attributes) == 0 {

		// Store cache
		logrus.WithField("collection", req.Collection).Infof("get data")
		data, err := getData(c, req.Collection)
		if err != nil {
			return fmt.Errorf("get data error: %v", err)
		}

		for _, item := range data.Results.AvailableAttributes {
			attributes = append(attributes, models.Attribute{
				Symbol:      data.Results.Symbol,
				TraitType:   item.Attribute.TraitType,
				Value:       item.Attribute.Value,
				Count:       item.Count,
				DateUpdated: time.Now().UTC(),
			})
		}

		logrus.WithField("collection", req.Collection).Infof("store")
		query := models.DB.NewInsert().Model(&attributes).On("CONFLICT(symbol, trait_type, value) DO UPDATE").Set("count = EXCLUDED.count")
		if _, err := query.Exec(ctx); err != nil {
			return fmt.Errorf("error query exec: %v", err)
		}

	}

	return c.JSON(http.StatusOK, attributes)
}

type Data struct {
	Results struct {
		Symbol              string `json:"symbol"`
		AvailableAttributes []struct {
			Count     int   `json:"count"`
			Floor     int64 `json:"floor"`
			Attribute struct {
				TraitType string `json:"trait_type"`
				Value     string `json:"value"`
			} `json:"attribute"`
		} `json:"availableAttributes"`
	} `json:"results"`
}

func getData(c echo.Context, collection string) (Data, error) {
	data := Data{}
	ctx := c.Request().Context()

	client, err := pkg.ByPassCloudFlareClient()
	if err != nil {
		return data, fmt.Errorf("error create cf bypass client: %v", err)
	}

	reqUrl := fmt.Sprintf("https://api-mainnet.magiceden.io/rpc/getCollectionEscrowStats/%s?status=all&edge_cache=true", collection)

	// Read cf pass cookie
	cfPass := schemas.AuthCookie{}
	cfPassCookie, err := models.Cache.Get(ctx, "cf_cookie").Bytes()
	if err != nil {
		return data, fmt.Errorf("get cf pass cookie error: %v", err)
	} else if err == redis.Nil || len(cfPassCookie) == 0 {
		return data, fmt.Errorf("not found cf pass cookie")
	} else if err := json.Unmarshal(cfPassCookie, &cfPass); err != nil {
		return data, fmt.Errorf("error decode redis key: %v", err)
	}

	r, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return data, fmt.Errorf("buy_now request create error: %v", err)
	}
	r.AddCookie(&http.Cookie{
		Name:     "__cf_bm",
		Value:    cfPass.Cookie,
		Path:     "/",
		Domain:   ".magiceden.io",
		Secure:   true,
		HttpOnly: true,
	})
	r.Header.Set("sec-ch-ua", cfPass.SecShUa)
	r.Header.Set("user-agent", cfPass.UserAgent)
	rr, err := client.Do(r)
	if err != nil {
		return data, fmt.Errorf("buy_now request error: %v", err)
	}
	defer rr.Body.Close()
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		return data, fmt.Errorf("buy_now request read error: %v", err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return data, fmt.Errorf("json unmarshal error: %v", err)
	}

	return data, nil
}
