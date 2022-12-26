package schemas

type CollectionsListRequest struct {
	Search          string `query:"search"`
	Symbol          string `query:"symbol" validate:"required"`
	MetaSymbol      string `query:"metaSymbol" validate:"required"`
	UpdateAuthority string `query:"updateAuthority" validate:"required"`
}
