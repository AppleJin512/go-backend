package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthCookie(c echo.Context) error {
	ctx := c.Request().Context()
	req := schemas.AuthCookie{}
	header := c.Request().Header.Get("Authorization")
	if header != "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhY2Nlc3MiLCJleHAiOjE2NTk5ODEyNTYsImlkIjoxLCJ1c2VybmFtZSI6IjEyMyJ9.TZoirjiIRNYjUEo0YdX_iKSRemlkmWmSxfQT6Ap0Cwk" {
		return c.JSON(http.StatusNotFound, schemas.Error{
			Message: "Error page not found",
		})
	}
	if err := c.Bind(&req); err != nil {
		return err
	} else if err := c.Validate(&req); err != nil {
		return err
	}

	if ok, err := CheckCookie(req); !ok {
		return c.JSON(http.StatusBadRequest, schemas.Error{
			Message: err.Error(),
		})
	}
	data, _ := json.Marshal(req)
	if err := models.Cache.Set(ctx, "cf_cookie", data, 0).Err(); err != nil {
		return fmt.Errorf("error set cache: %v", err)
	}
	return c.JSON(http.StatusOK, req)
}

func CheckCookie(r schemas.AuthCookie) (bool, error) {

	client, err := pkg.ByPassCloudFlareClient()
	if err != nil {
		return false, fmt.Errorf("error create cf bypass client: %v", err)
	}

	url := "https://api-mainnet.magiceden.io/rpc/getCollectionEscrowStats/sniper_cove?status=all&edge_cache=true"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("error create request: %v", err)
	}
	req.AddCookie(&http.Cookie{
		Name:     "__cf_bm",
		Value:    r.Cookie,
		Path:     "/",
		Domain:   ".magiceden.io",
		Secure:   true,
		HttpOnly: true,
	})
	req.Header.Set("sec-ch-ua", r.SecShUa)
	req.Header.Set("user-agent", r.UserAgent)
	fmt.Printf("%s\n", r.Cookie)
	fmt.Printf("%s\n", r.SecShUa)
	fmt.Printf("%s\n", r.UserAgent)
	response, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error get do: %v", err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("read error: %v", err)
	}

	resp := Series{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return false, fmt.Errorf("decode error: %v", err)
	}

	return resp.Results.Symbol == "sniper_cove", nil
}

type Series struct {
	Results struct {
		Symbol string `json:"symbol"`
	} `json:"results"`
}
