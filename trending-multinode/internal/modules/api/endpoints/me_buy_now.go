package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"
	"net/url"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

func MEBuyNow(c echo.Context) error {
	ctx := c.Request().Context()
	req := schemas.MEBuyNowRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	} else if err := c.Validate(&req); err != nil {
		return err
	}
	client, err := pkg.ByPassCloudFlareClient()
	if err != nil {
		return fmt.Errorf("error create cf bypass client: %v", err)
	}

	// @TODO maybe use orderered params
	params := url.Values{}
	params.Add("buyer", req.Buyer)
	params.Add("seller", req.Seller)
	params.Add("auctionHouseAddress", req.AuctionHouseAddress)
	params.Add("tokenMint", req.TokenMint)
	params.Add("tokenATA", req.TokenATA)
	params.Add("price", req.Price)
	params.Add("sellerReferral", req.SellerReferral)
	params.Add("sellerExpiry", "-1")
	reqUrl := "https://api-mainnet.magiceden.io/v2/instructions/buy_now" + "?" + params.Encode()

	// Read cf pass cookie
	cfPass := schemas.AuthCookie{}
	cfPassCookie, err := models.Cache.Get(ctx, "cf_cookie").Bytes()
	if err != nil {
		return fmt.Errorf("get cf pass cookie error: %v", err)
	} else if err == redis.Nil || len(cfPassCookie) == 0 {
		return fmt.Errorf("not found cf pass cookie")
	} else if err := json.Unmarshal(cfPassCookie, &cfPass); err != nil {
		return fmt.Errorf("error decode redis key: %v", err)
	}

	r, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return fmt.Errorf("buy_now request create error: %v", err)
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
		return fmt.Errorf("buy_now request error: %v", err)
	}
	defer rr.Body.Close()
	data, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		return fmt.Errorf("buy_now request read error: %v", err)
	}

	return c.Blob(http.StatusOK, "application/json", data)
}
