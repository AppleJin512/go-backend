package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Attribute struct {
	bun.BaseModel `bun:"table:attributes,alias:a"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id" validate:"-"`
	Symbol      string    `json:"symbol"`
	TraitType   string    `json:"traitType"`
	Value       string    `json:"value"`
	Count       int       `json:"count"`
	DateUpdated time.Time `json:"dateUpdated"`
}
