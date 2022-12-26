package schemas

type MEBuyNowRequest struct {
	Buyer               string `query:"buyer"`
	Seller              string `query:"seller"`
	TokenMint           string `query:"token_mint"`
	TokenATA            string `query:"token_ata"`
	Price               string `query:"price"`
	SellerReferral      string `query:"seller_referral"`
	AuctionHouseAddress string `query:"auction_house_address"`
}
