package schemas

import "moonbite/trending/internal/models"

type ListingsListRequest struct {
	Symbol string `query:"symbol"`
	Offset int    `query:"offset"`
	Limit  int    `query:"limit"`
	Order  string `query:"order"`
}

type ListingsListResponse struct {
	Items []models.Listing `json:"items"`
}
