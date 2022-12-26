package pkg

import (
	"moonbite/trending/internal/config"
	"net/http"
	"net/url"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/memclutter/gorequests"
	"github.com/sirupsen/logrus"
)

func GetProxies() []*url.URL {
	proxies := make([]*url.URL, 0)
	for _, p := range config.Proxies {
		proxyUrl, err := url.Parse(p)
		if err != nil {
			logrus.Warnf("invalid proxy %s: %v", p, err)
			continue
		}
		proxies = append(proxies, proxyUrl)
	}
	return proxies
}

func ByPassCloudFlareClient() (*http.Client, error) {
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

	proxyUrl, _ := url.Parse("socks5://R1EiaETT:ysrb2Ekf@154.3.111.206:56179")

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport, options)

	return client, nil
}

func CheckProxy(pr *url.URL) error {
	return gorequests.Get(
		gorequests.WithExtensions(
			gorequests.ProxiesExtension{Proxies: []string{pr.String()}},
		),
		gorequests.WithUrl("https://trending.moonbite.io/collections/?symbol=okay_bears"),
		gorequests.WithOkStatusCodes(http.StatusOK),
	)
}
