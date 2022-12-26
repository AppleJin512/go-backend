package schemas

import "moonbite/trending/internal/models"

type ActivitiesListRequest struct {
	Symbol       string `query:"symbol"`
	ActivityType string `query:"activity_type"`
	Offset       int    `query:"offset"`
	Limit        int    `query:"limit"`
}

type ActivitiesListResponse struct {
	Items []models.Activity `json:"items"`
}
