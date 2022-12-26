package models

import (
	"encoding/json"
	"github.com/uptrace/bun"
	"time"
)

type Listing struct {
	bun.BaseModel `bun:"table:listings,alias:l"`

	ID                  int64           `bun:"id,pk,autoincrement" json:"id" validate:"-"`
	Signature           string          `json:"signature"`
	Symbol              string          `json:"symbol"`
	BlockTime           time.Time       `json:"block_time"`
	Price               float64         `json:"price"`
	Name                string          `json:"name"`
	MintAddress         string          `json:"mint_address"`
	Uri                 string          `json:"uri"`
	Seller              string          `json:"seller"`
	AuctionHouseAddress string          `json:"auctionHouseAddress"`
	SellerReferral      string          `json:"sellerReferral"`
	Rank                int             `json:"rank"`
	Attributes          json.RawMessage `json:"attributes"`
}
