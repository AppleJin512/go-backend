package blockchain

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	JobTypeSub string = "SUB" // log subscribe
	JobTypeLtp string = "LTP" // long time pulling
)

const (
	TrxTypeUnknown string = "TrxTypeUnknown"

	// MagicEden
	TrxTypeMeListingV2            string = "TrxTypeMeListingV2"
	TrxTypeMeListingUpdatePriceV2 string = "TrxTypeMeListingUpdatePriceV2"
	TrxTypeMeBuyV2                string = "TrxTypeMeBuyV2"
	TrxTypeMeDelistingV2          string = "TrxTypeMeDelistingV2"

	// Solanart

	TrxTypeSaListing     string = "TrxTypeSaListing"
	TrxTypeSaUpdatePrice string = "TrxTypeSaUpdatePrice"
)

const (
	QuickNode = "https://weathered-broken-bush.solana-mainnet.quiknode.pro/dc821bd2b1663e07d4d95257c92b797379385d13/"

	GalaxityNode   = "https://rpc.galaxity.io/"
	GalaxityNodeWS = "wss://rpc.galaxity.io/"

	Xenode   = "http://47.253.131.10:8899/"
	XenodeWS = "ws://47.253.131.10:8900"

	P9   = "http://b1be1602695a1d56ec206258aa2c466d.p9nodes.io:8899/"
	P9WS = "ws://b1be1602695a1d56ec206258aa2c466d.p9nodes.io:8900/"

	NodeMonkey   = "http://genesysnode1.nodemonkey.io:8899/"
	NodeMonkeyWS = "ws://genesysnode1.nodemonkey.io:8900/"
)

var (
	MESmartContractV1 = solana.MustPublicKeyFromBase58("MEisE1HzehtrDpAAT8PnLHjpSSkRYakotTuJRPjTpo8")
	MESmartContractV2 = solana.MustPublicKeyFromBase58("M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K")
	SASmartContract   = solana.MustPublicKeyFromBase58("CJsLwbP1iu5DuUikHEJnLfANgKy6stB2uFgvBBHoyxwz")
	FRSmartContract   = solana.MustPublicKeyFromBase58("hausS13jsjafwWwGqZTUQRmWyvyxn9EQpqMwV1PBBmk")

	TokenMetadataProgramID = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
)

type NodeSetting struct {
	Contract   solana.PublicKey
	JobType    string
	NodeUrlRpc string
	NodeUrlWs  string
}

var NodeSettings = map[string]NodeSetting{
	"ME_V2_LTP_SERUM": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeLtp,
		NodeUrlRpc: rpc.MainNetBetaSerum_RPC,
		NodeUrlWs:  rpc.MainNetBetaSerum_WS,
	},
	"ME_V2_SUB_SERUM": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: rpc.MainNetBetaSerum_RPC,
		NodeUrlWs:  rpc.MainNetBetaSerum_WS,
	},
	"ME_V2_SUB_GENESIS": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "https://ssc-dao.genesysgo.net/",
		NodeUrlWs:  "wss://ssc-dao.genesysgo.net/",
	},
	"GALAXITY_SUB": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "https://rpc.galaxity.io",
		NodeUrlWs:  "wss://rpc.galaxity.io",
	},
	"GALAXITY_EU_SUB": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "https://rpc-3.galaxity.io",
		NodeUrlWs:  "wss://rpc-3.galaxity.io",
	},
	"XENODES_SUB": {
		Contract: MESmartContractV2,
		JobType:  JobTypeSub,
		//NodeUrlRpc: "https://xenodes-andresupbeat.xenonx.io/",
		NodeUrlRpc: "https://47.253.131.10:8899",
		//NodeUrlWs:  "ws://xenodes-andresupbeat.xenonx.io",
		NodeUrlWs: "ws://47.253.131.10:8900",
	},
	"GEN_NODE": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "http://genesysnode1.nodemonkey.io:8899/",
		NodeUrlWs:  "ws://genesysnode1.nodemonkey.io:8900/",
	},
	"GENESIS_NODE": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "https://ssc-dao.genesysgo.net/",
		NodeUrlWs:  "wss://ssc-dao.genesysgo.net/",
	},
	"P9_NODE": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "http://b1be1602695a1d56ec206258aa2c466d.p9nodes.io:8899/",
		NodeUrlWs:  "ws://b1be1602695a1d56ec206258aa2c466d.p9nodes.io:8900/",
	},
	"XENONX_SUB": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "https://ny-testnode.xenonx.io",
		NodeUrlWs:  "ws://ny-testnode.xenonx.io:8900/",
	},
	"XENONX_LTP": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeLtp,
		NodeUrlRpc: "https://ny-testnode.xenonx.io",
		NodeUrlWs:  "ws://ny-testnode.xenonx.io:8900/",
	},
	"ALCHEMY_SUB": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "https://solana-mainnet.g.alchemy.com/v2/zoF5qa7Xs-XD5pMe4NOSfl2QEA61phD2",
		NodeUrlWs:  "wss://solana-mainnet.g.alchemy.com/v2/zoF5qa7Xs-XD5pMe4NOSfl2QEA61phD2",
	},
	"ALCHEMY_LTP": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeLtp,
		NodeUrlRpc: "https://solana-mainnet.g.alchemy.com/v2/zoF5qa7Xs-XD5pMe4NOSfl2QEA61phD2",
		NodeUrlWs:  "wss://solana-mainnet.g.alchemy.com/v2/zoF5qa7Xs-XD5pMe4NOSfl2QEA61phD2",
	},
	"ZONERNODES_SUB": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeSub,
		NodeUrlRpc: "https://rpc-va.zonernodes.io/",
		NodeUrlWs:  "wss://rpc-va.zonernodes.io",
	},
	"ZONERNODES_LTP": {
		Contract:   MESmartContractV2,
		JobType:    JobTypeLtp,
		NodeUrlRpc: "https://rpc-va.zonernodes.io/",
		NodeUrlWs:  "wss://rpc-va.zonernodes.io",
	},
}
