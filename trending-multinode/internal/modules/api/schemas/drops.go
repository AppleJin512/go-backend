package schemas

type DropsRequest struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}
