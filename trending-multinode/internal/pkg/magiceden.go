package pkg

import (
	"context"
	"fmt"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/models"
	"net/http"
	"time"

	bin "github.com/gagliardetto/binary"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/memclutter/gorequests"
	"github.com/sirupsen/logrus"
)

type MagicResponse struct {
	TokenMint string `json:"tokenMint"`
}

func GetCollectionData(ctx context.Context, rpcClient *rpc.Client, symbol string, collection *models.Collection) error {

	// get update authority
	listings := make([]MagicResponse, 0)
	if err := gorequests.Get(
		gorequests.WithExtensions(
			gorequests.ProxiesExtension{Proxies: config.Proxies},
			gorequests.RetryExtension{RetryMax: 5, RetryWaitMin: 500 * time.Millisecond, RetryWaitMax: 800 * time.Millisecond},
		),
		gorequests.WithUrl("https://api-mainnet.magiceden.dev/v2/collections/%s/listings?offset=0&limit=1", symbol),
		gorequests.WithOkStatusCodes(http.StatusOK),
		gorequests.WithOut(&listings, gorequests.OutTypeJson),
	); err != nil {
		return fmt.Errorf("gorequests get error: %v", err)
	}

	if len(listings) == 0 {
		// Get activities
		time.Sleep(500 * time.Millisecond)
		if err := gorequests.Get(
			gorequests.WithExtensions(
				gorequests.ProxiesExtension{Proxies: config.Proxies},
				gorequests.RetryExtension{RetryMax: 5, RetryWaitMin: 500 * time.Millisecond, RetryWaitMax: 800 * time.Millisecond},
			),
			gorequests.WithUrl("https://api-mainnet.magiceden.dev/v2/collections/%s/activities?offset=0&limit=1", symbol),
			gorequests.WithOkStatusCodes(http.StatusOK),
			gorequests.WithOut(&listings, gorequests.OutTypeJson),
		); err != nil {
			return fmt.Errorf("gorequests get error: %v", err)
		}

		// Not found
		if len(listings) == 0 {

			// https://api-mainnet.magiceden.io/rpc/getListedNFTsByQueryLite?q=%7B%22%24match%22%3A%7B%22%24or%22%3A%5B%7B%22collectionSymbol%22%3A%22benny_andallo_fuzz%22%7D%2C%7B%22onChainCollectionAddress%22%3A%225KNR4isxnzUk6igkRHoWSjYAyYckKsPbgoezED5xLJR6%22%7D%5D%7D%2C%22%24sort%22%3A%7B%22takerAmount%22%3A1%7D%2C%22%24skip%22%3A0%2C%22%24limit%22%3A1%2C%22status%22%3A%5B%22all%22%5D%7D

			logrus.Warnf("empty update authority %s", symbol)
			return nil
		}
	}

	mintAddress := listings[0].TokenMint
	metadataAddress := GetSolanaMetadataAddress(solana.MustPublicKeyFromBase58(mintAddress))

	opts := &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
		Encoding:   "base64",
	}

	r, err := rpcClient.GetAccountInfoWithOpts(ctx, metadataAddress, opts)
	if err != nil {
		return fmt.Errorf("get account info error: %v", err)
	}

	data := r.Value.Data.GetBinary()
	dec := bin.NewBorshDecoder(data)

	var meta token_metadata.Metadata
	err = dec.Decode(&meta)
	if err != nil {
		return fmt.Errorf("decode account info error: %v", err)
	}

	collection.MetaSymbol = ClearInvisibleChars(meta.Data.Symbol)
	collection.UpdateAuthority = meta.UpdateAuthority.String()

	return nil
}
