package schemas

type MESeriesRequest struct {
	Period     string `json:"period" query:"period"`
	Collection string `json:"collection" query:"collection"`
}
