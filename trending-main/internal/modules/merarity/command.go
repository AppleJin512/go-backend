package merarity

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/modules/api/schemas"
	"moonbite/trending/internal/pkg"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Command(c *cli.Context) error {
	ctx := c.Context
	for {
		if err := Run(ctx); err != nil {
			logrus.Errorf("error run: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}
		time.Sleep(10 * time.Minute)
	}

}

func Run(ctx context.Context) error {

	logrus.Info("run")

	interval := time.Now().UTC().Add(-30 * 1 * time.Hour) // 30 days???
	collections := make([]models.Collection, 0)
	offset := 0
	limit := 10000
	query := models.DB.NewSelect().Model(&collections).OrderExpr("date_last_items_update DESC").Where("date_last_items_update > ? OR date_last_items_update IS NULL", interval)

	for {
		if err := query.Offset(offset).Limit(limit).Scan(ctx); err == sql.ErrNoRows {
			break
		} else if err != nil {
			return err
		}

		if len(collections) == 0 {
			break
		}

		for _, collection := range collections {
			if err := RunCollection(ctx, &collection); err != nil {
				logrus.Errorf("run collections error: %v", err)
			} else {
				collection.DateLastItemsUpdate = sql.NullTime{
					Time:  time.Now().UTC(),
					Valid: true,
				}
				if _, err := models.DB.NewUpdate().Model(&collection).WherePK().Exec(ctx); err != nil {
					logrus.Errorf("save collection error: %v", err)
				}
			}
		}

		offset += limit
	}

	return nil
}

func RunCollection(ctx context.Context, collection *models.Collection) error {

	logrus.Infof("run collections %s", collection.Symbol)
	// https://api-mainnet.magiceden.dev/v2/collections/degenerate_ape_academy/listings?offset=0&limit=20
	offset := 0
	limit := 20

	retries := 10

	if _, err := models.DB.NewDelete().Model((*models.Item)(nil)).Where("collection_symbol = ?", collection.Symbol).Exec(ctx); err != nil {
		return fmt.Errorf("delete error: %v", err)
	}

	for {
		result, err := GetListedNFT(ctx, collection.Symbol, offset, limit)
		if err != nil {
			logrus.Errorf("error get listed nt: %v", err)
			if retries <= 0 {
				return fmt.Errorf("max retries exceed: %v", err)
			}
			retries -= 1
			time.Sleep(30 * time.Second)
			continue
		}

		if len(result.Results) == 0 {
			break
		}

		// Parse items
		items := make([]models.Item, 0)
		for _, re := range result.Results {
			attributes, _ := json.Marshal(re.Attributes)
			items = append(items, models.Item{
				CollectionSymbol: collection.Symbol,
				TokenMint:        re.MintAddress,
				Title:            re.Title,
				Img:              re.Img,
				Rank:             re.Rarity.Merarity.Rank,
				Attributes:       attributes,
			})
		}

		if len(items) > 0 {
			logrus.WithField("collection", collection.Symbol).Infof("insert items %d", len(items))
			if _, err := models.DB.NewInsert().Model(&items).Exec(ctx); err != nil {
				return fmt.Errorf("insert error: %v", err)
			}
		}

		offset += limit
	}

	return nil
}

type ListedNFTResultAttribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

type ListedNFTResultRarity struct {
	Merarity struct {
		TokenKey    string  `json:"tokenKey"`
		Score       float64 `json:"score"`
		TotalSupply int     `json:"totalSupply"`
		Rank        int     `json:"rank"`
	} `json:"merarity"`
}

type PreListedNFTResult struct {
	MintAddress string                `json:"mintAddress"`
	Attributes  []json.RawMessage     `json:"attributes"`
	CreatedAt   time.Time             `json:"createdAt"`
	Img         string                `json:"img"`
	Title       string                `json:"title"`
	Rarity      ListedNFTResultRarity `json:"rarity"`
}

type ListedNFTResult struct {
	MintAddress string                     `json:"mintAddress"`
	Attributes  []ListedNFTResultAttribute `json:"attributes"`
	CreatedAt   time.Time                  `json:"createdAt"`
	Img         string                     `json:"img"`
	Title       string                     `json:"title"`
	Rarity      ListedNFTResultRarity      `json:"rarity"`
}

type PreListedNFT struct {
	Results []PreListedNFTResult `json:"results"`
}

type ListedNFT struct {
	Results []ListedNFTResult `json:"results"`
}

func GetListedNFT(ctx context.Context, collection string, offset, limit int) (ListedNFT, error) {
	result := ListedNFT{}
	resultPre := PreListedNFT{}

	// Format url
	q := fmt.Sprintf(`{"$match":{"collectionSymbol":"%s","rarity.merarity":{"$exists":true}},"$sort":{"rarity.merarity.rank":-1},"$skip":%d,"$limit":%d,"status":["all"]}`, collection, offset, limit)
	p := url.Values{}
	p.Set("q", q)
	u := "https://api-mainnet.magiceden.io/rpc/getListedNFTsByQueryLite?" + p.Encode()
	fmt.Println(u)

	cl, err := pkg.ByPassCloudFlareClient()
	if err != nil {
		return result, fmt.Errorf("error create cf bypass client: %v", err)
	}

	// Read cf pass cookie
	logrus.WithField("collection", collection).Info("get cookie")
	cfPass := schemas.AuthCookie{}
	cfPassCookie, err := models.Cache.Get(ctx, "cf_cookie").Bytes()
	if err != nil {
		return result, fmt.Errorf("get cf pass cookie error: %v", err)
	} else if err == redis.Nil || len(cfPassCookie) == 0 {
		return result, fmt.Errorf("not found cf pass cookie")
	} else if err := json.Unmarshal(cfPassCookie, &cfPass); err != nil {
		return result, fmt.Errorf("error decode redis key: %v", err)
	}

	logrus.WithField("collection", collection).Info("request")
	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return result, fmt.Errorf("buy_now request create error: %v", err)
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
	rr, err := cl.Do(r)
	if err != nil {
		return result, fmt.Errorf("buy_now request error: %v", err)
	}
	defer rr.Body.Close()
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		return result, fmt.Errorf("buy_now request read error: %v", err)
	}

	if err := json.Unmarshal(body, &resultPre); err != nil {
		return result, fmt.Errorf("json unmarshal error: %v, body is %s", err, body)
	}

	result.Results = make([]ListedNFTResult, 0)
	for _, r := range resultPre.Results {
		attributes := make([]ListedNFTResultAttribute, 0)
		for _, a := range r.Attributes {
			aa := ListedNFTResultAttribute{}
			if err := json.Unmarshal(a, &aa); err == nil {
				attributes = append(attributes, aa)
			}
		}

		result.Results = append(result.Results, ListedNFTResult{
			MintAddress: r.MintAddress,
			Attributes:  attributes,
			CreatedAt:   r.CreatedAt,
			Img:         r.Img,
			Title:       r.Title,
			Rarity:      r.Rarity,
		})
	}

	return result, nil
}
