package sol

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"net/url"
	"strings"
)

type NodeEnv struct {
	Url            string
	Mentions       solana.PublicKey
	CommitmentType rpc.CommitmentType
}

var (
	DefaultMentions       = solana.MustPublicKeyFromBase58("M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K")
	DefaultCommitmentType = rpc.CommitmentConfirmed
)

// ParseNodeEnvSlice godoc
func ParseNodeEnvSlice(slice []string) (resSlice []NodeEnv, err error) {
	for i, val := range slice {
		if res, err := ParseNodeEnv(val); err != nil {
			return resSlice, fmt.Errorf("[%d]: %v", i, err)
		} else {
			resSlice = append(resSlice, res)
		}
	}
	return
}

// ParseNodeEnv godoc
//
// The function parses a string like
// - (https?|wss?)://solana.node.example.com/?mentions=(public_key)&commitmentType=confirmed
// and gets the value of NodeEnv. The final result is used to create connections to Solana nodes.
func ParseNodeEnv(val string) (res NodeEnv, err error) {
	u, err := url.Parse(val)
	if err != nil {
		return res, fmt.Errorf("invalid url: %v", err)
	}
	res.Mentions = DefaultMentions
	res.CommitmentType = DefaultCommitmentType
	q := u.Query()
	if mentions := q.Get("mentions"); len(mentions) > 0 {
		if pk, err := solana.PublicKeyFromBase58(mentions); err != nil {
			return res, fmt.Errorf("invalid mentions public '%s' key: %v", mentions, err)
		} else {
			res.Mentions = pk
		}
	}
	if commitmentTypeStr := q.Get("commitmentType"); len(commitmentTypeStr) > 0 {
		commitmentType := rpc.CommitmentType(commitmentTypeStr)
		if commitmentType != rpc.CommitmentConfirmed && commitmentType != rpc.CommitmentFinalized {
			return res, fmt.Errorf("unsupported commitment type '%s', available 'confirmed', 'finalized'", commitmentType)
		} else {
			res.CommitmentType = commitmentType
		}
	}
	q.Del("mentions")
	q.Del("commitmentType")
	u.ForceQuery = true
	u, _ = url.Parse(strings.Split(u.String(), "?")[0] + "?" + q.Encode())
	u.ForceQuery = false
	res.Url = u.String()
	return
}
