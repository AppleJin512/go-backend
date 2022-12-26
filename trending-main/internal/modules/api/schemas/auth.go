package schemas

type AuthCookie struct {
	Cookie    string `query:"cookie" json:"cookie"`
	UserAgent string `query:"userAgent" json:"userAgent"`
	SecShUa   string `query:"secShUa" json:"secShUa"`
}
