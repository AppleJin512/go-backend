package schemas

import (
	"github.com/shopspring/decimal"
	"time"
)

type StatsRequest struct {
	Period string `json:"period" query:"period"`
	Sort   string `json:"sort" query:"sort"`
}

type StatsResponse struct {
	Symbol            string          `json:"symbol"`
	Name              string          `json:"name"`
	Image             string          `json:"image"`
	StatsVolume       decimal.Decimal `json:"stats_volume"`
	StatsTradeCount   int             `json:"stats_trade_count"`
	StatsAveragePrice decimal.Decimal `json:"stats_average_price"`
}

type StatsDistributionRequest struct {
	Period string `query:"period"`
	Symbol string `query:"symbol"`
	Limit  int    `query:"limit"`
}

type StatsDistributionResponse struct {
	Slot         time.Time `json:"slot"`
	CountSale    int64     `json:"count_sale"`
	CountListing int64     `json:"count_listing"`
}
