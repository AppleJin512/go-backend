package sol

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/test-go/testify/assert"
	"testing"
)

func TestParseNodeEnv(t *testing.T) {
	cases := []struct {
		name string
		val  string
		res  NodeEnv
	}{
		{
			name: "can parse https",
			val:  "https://ssc-dao.genesysgo.net/?mentions=M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K&commitmentType=confirmed",
			res: NodeEnv{
				Url:            "https://ssc-dao.genesysgo.net/",
				Mentions:       solana.MustPublicKeyFromBase58("M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K"),
				CommitmentType: rpc.CommitmentConfirmed,
			},
		},
		{
			name: "can parse wss",
			val:  "wss://solana-mainnet.g.alchemy.com/v2/zoF5qa7Xs-XD5pMe4NOSfl2QEA61phD2",
			res: NodeEnv{
				Url:            "wss://solana-mainnet.g.alchemy.com/v2/zoF5qa7Xs-XD5pMe4NOSfl2QEA61phD2",
				Mentions:       DefaultMentions,
				CommitmentType: DefaultCommitmentType,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ParseNodeEnv(tc.val)
			assert.NoError(t, err, "should no errors")
			assert.Equal(t, tc.res, res, "should be equal result")
		})
	}
}
