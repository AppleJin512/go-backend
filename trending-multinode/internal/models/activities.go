package models

import (
	"moonbite/trending/internal/blockchain"
	"time"

	"github.com/uptrace/bun"
)

type Activity struct {
	bun.BaseModel `bun:"table:activities,alias:a"`

	ID                  int64     `bun:"id,pk,autoincrement" json:"id" validate:"-"`
	Signature           string    `json:"signature"`
	Symbol              string    `json:"symbol"`
	BlockTime           time.Time `json:"block_time"`
	Price               float64   `json:"price"`
	ActivityType        string    `json:"activity_type"`
	Name                string    `json:"name"`
	MintAddress         string    `json:"mint_address"`
	Uri                 string    `json:"uri"`
	Seller              string    `json:"seller"`
	AuctionHouseAddress string    `json:"auctionHouseAddress"`
	SellerReferral      string    `json:"sellerReferral"`

	Item *Item `json:"item" bun:"rel:belongs-to,join:mint_address=token_mint"`
}

const (
	ActivityTypeDelisting   = "delisting"
	ActivityTypeListing     = "listing"
	ActivityTypeSale        = "sale"
	ActivityTypeUpdatePrice = "update_price"
)

var TrxTypeMapping = map[string]string{
	blockchain.TrxTypeMeDelistingV2:          ActivityTypeDelisting,
	blockchain.TrxTypeMeListingV2:            ActivityTypeListing,
	blockchain.TrxTypeMeBuyV2:                ActivityTypeSale,
	blockchain.TrxTypeMeListingUpdatePriceV2: ActivityTypeUpdatePrice,
}
