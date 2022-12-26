package models

import (
	"database/sql"
	"github.com/uptrace/bun"
)

type Collection struct {
	bun.BaseModel `bun:"table:collections,alias:c"`

	ID                  int64        `bun:"id,pk,autoincrement" json:"id" validate:"-"`
	Symbol              string       `json:"symbol"`
	Name                string       `json:"name"`
	Image               string       `json:"image"`
	MetaSymbol          string       `json:"metaSymbol"`
	UpdateAuthority     string       `json:"updateAuthority"`
	DateLastItemsUpdate sql.NullTime `json:"-"`
}
