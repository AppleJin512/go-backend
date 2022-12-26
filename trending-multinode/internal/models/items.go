package models

import (
	"encoding/json"
	"github.com/uptrace/bun"
)

type Item struct {
	bun.BaseModel `bun:"table:items,alias:i"`

	ID               int64           `bun:"id,pk,autoincrement" json:"id" validate:"-"`
	CollectionSymbol string          `json:"collection_symbol"`
	TokenMint        string          `json:"token_mint"`
	Title            string          `json:"title"`
	Img              string          `json:"img"`
	Rank             int             `json:"rank"`
	Attributes       json.RawMessage `json:"attributes"`
}
