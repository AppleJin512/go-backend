package stats

import "github.com/shopspring/decimal"

type Activity struct {
	Signature        string          `json:"signature"`
	CollectionSymbol string          `json:"collectionSymbol"`
	BlockTime        int             `json:"blockTime"`
	Price            decimal.Decimal `json:"price"`
}
