package cookies

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/endpoints"
	"moonbite/trending/internal/modules/api/schemas"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/docker/docker/api/types"
	client "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Command(c *cli.Context) error {
	for {
		if err := FetchCfCookie(context.Background()); err != nil {
			logrus.Errorf("fetch cookie error: %v", err)
			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(20 * time.Minute)
		}
	}
}

func FetchCfCookie(ctx context.Context) error {
	logrus.Info("fetch cookie")
	d, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.39"))
	if err != nil {
		return fmt.Errorf("error connect docker client: %v", err)
	}

	logsReader, err := d.ContainerLogs(ctx, "chrome", types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      "",
		Until:      "",
		Timestamps: false,
		Follow:     false,
		Tail:       "200",
		Details:    false,
	})
	if err != nil {
		return fmt.Errorf("error get logs: %v", err)
	}
	defer logsReader.Close()
	logs, _ := ioutil.ReadAll(logsReader)

	re := regexp.MustCompile(`(ws:.*)`)
	urls := re.FindAll(logs, -1)
	if len(urls) == 0 {
		return fmt.Errorf("not found cdp url")
	}
	cdpUrl := strings.ReplaceAll(string(urls[len(urls)-1]), "127.0.0.1:34001", "chrome:34000")

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), cdpUrl)
	defer cancel()

	ctxWithLogs, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	var res string
	err = chromedp.Run(ctxWithLogs,
		chromedp.Navigate("https://magiceden.io/launchpad/about"),
		chromedp.Sleep(2*time.Second),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}

			for _, cookie := range cookies {
				if cookie.Name == "__cf_bm" {
					logrus.Infof("found cf cookie in browser: %s", cookie.Value)
					res = cookie.Value
				}
			}

			return nil
		}),
		chromedp.Stop(),
	)
	if err != nil {
		return fmt.Errorf("error run cdp script: %v", err)
	}

	req := schemas.AuthCookie{
		Cookie:    res,
		SecShUa:   `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`,
		UserAgent: `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36`,
	}
	if ok, err := endpoints.CheckCookie(req); !ok {
		return fmt.Errorf("error check cookie: %v", err)
	}
	data, _ := json.Marshal(req)
	if err := models.Cache.Set(ctx, "cf_cookie", data, 0).Err(); err != nil {
		return fmt.Errorf("error set cache: %v", err)
	}
	logrus.Infof("set cf cookie in redis %s", req.Cookie)
	return nil
}
