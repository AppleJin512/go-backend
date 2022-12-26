package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"moonbite/trending/internal/config"
	"net/http"
	"net/url"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

func main() {

	options := cloudflarebp.Options{
		AddMissingHeaders: true,
		Headers: map[string]string{
			"authority":          "api-mainnet.magiceden.io",
			"method":             "GET",
			"path":               "/v2/instructions/buy_now?buyer=BYtHXofFczcvcsbLMRQHwAzKQMqzg393kV8tuHExhpd2&seller=AWiso2JYQ8SLR8Pe3zVsMqq3hUgkzXpC3uRsatpGFJY9&auctionHouseAddress=E8cU1WiRWjanGxmn96ewBgk9vPTcL6AEZ1t6F6fkgUWe&tokenMint=HZZ8noZLMT59Axu3JZFBiX2yCzbpdihp7Tk9TxwCPBmM&tokenATA=EC1SZ6opmG1wiPNQehEGm81KCE5f9M3eoSubf93mJ8wu&price=0.675&sellerReferral=autMW8SgBkVYeBgqYiTuJZnkvDZMVU2MHJh9Jh7CSQ2&sellerExpiry=-1",
			"scheme":             "https",
			"accept":             "application/json, text/plain, */*",
			"accept-language":    "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			"if-none-match":      `W/"17e2-EN+a8jNdChzn6psCKlcwVCMfxPg"`,
			"origin":             "https://www.magiceden.io",
			"referer":            "https://www.magiceden.io/",
			"sec-ch-ua":          `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`,
			`sec-ch-ua-mobile`:   `?0`,
			`sec-ch-ua-platform`: `"Linux"`,
			`sec-fetch-dest`:     `empty`,
			`sec-fetch-mode`:     `cors`,
			`sec-fetch-site`:     `same-site`,
			`user-agent`:         `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36`,
		},
	}

	req, err := http.NewRequest(http.MethodGet, "https://stats-mainnet.magiceden.io/collection_stats/getCollectionTimeSeries/okay_bears?edge_cache=true&resolution=1h&addLastDatum=true", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.AddCookie(&http.Cookie{
		Name:     "__cf_bm",
		Value:    "cUoJEhCqWuqvv4_v5DJZQj4r.M_r92CDYPxH9UfC9B8-1666304162-0-Aciju6hBVqGzZMe1QQrIr2+IvcRRqJfeEqFxYXxDLjddbBP+OJrOVCv1PwYEfKgGzEAjknQuWc+zGBvcyUkf5fm5K6+ozEK5unvxSFXTZsQnVKyjhBhbDi7XA1+Px1rG4EjXY/jIwzS22k5vNDBWhcBwkqqDKcf3Oal0rkw5kWGY",
		Path:     "/",
		Domain:   ".magiceden.io",
		Secure:   true,
		HttpOnly: true,
	})

	for _, proxy := range config.Proxies {
		proxyUrl, _ := url.Parse(proxy)

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport, options)

		response, err := client.Do(req)
		if err != nil {
			fmt.Printf("%s: fail %v\n", proxy, err)
			continue
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s: fail %v\n", proxy, err)
			continue
		}

		//spew.Dump(data)

		resp := make([]Resp, 0)
		if err := json.Unmarshal(data, &resp); err != nil {
			fmt.Printf("%s: fail %v\n", proxy, err)
			continue
		}

		if len(resp) == 0 {
			fmt.Printf("%s: fail %s\n", proxy, data[0:40])
			continue
		}
		fmt.Printf("%s: ok\n", proxy)
	}

}

type Resp struct {
	CFP   float64 `json:"cFP"`
	CLC   float64 `json:"cLC"`
	CV    float64 `json:"cV"`
	MaxFP float64 `json:"maxFP"`
	MinFP float64 `json:"minFP"`
	OFP   float64 `json:"oFP"`
	OLC   float64 `json:"oLC"`
	OV    float64 `json:"oV"`
	Ts    float64 `json:"ts"`
}
