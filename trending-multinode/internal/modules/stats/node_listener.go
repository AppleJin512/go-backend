package stats

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"moonbite/trending/internal/blockchain"
	"moonbite/trending/internal/config"
	"moonbite/trending/internal/models"
	"moonbite/trending/internal/pkg"
	"net/http"
	"strings"
	"sync"
	"time"

	bin "github.com/gagliardetto/binary"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/sirupsen/logrus"
)

type NodeListener struct {
	nodeSetting blockchain.NodeSetting
	sub         *ws.LogSubscription
	rpcClient   *rpc.Client
	wsClient    *ws.Client

	Log *logrus.Entry

	signatureCache sync.Map
	ctx            context.Context
}

func NewNodeListeners(ctx context.Context, setting string) (*rpc.Client, []*NodeListener, error) {
	var operationRpcClient *rpc.Client
	listeners := make([]*NodeListener, 0)
	for i, name := range strings.Split(setting, ":") {
		ns, ok := blockchain.NodeSettings[name]
		if !ok {
			return nil, nil, fmt.Errorf("Can't parse setting: '%s'\n", name)
		}
		instance := &NodeListener{
			ctx:            ctx,
			nodeSetting:    ns,
			sub:            &ws.LogSubscription{},
			signatureCache: sync.Map{},
			Log: logrus.WithFields(logrus.Fields{
				"component": "node_listener",
				"node":      name,
				"job_type":  ns.JobType,
			}),
		}

		// Init node connects
		instance.rpcClient = rpc.New(ns.NodeUrlRpc)
		if i == 0 {
			operationRpcClient = instance.rpcClient
		}
		_, err := instance.rpcClient.GetHealth(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("error get health of node %s: %v", name, err)
		}
		if ns.JobType == blockchain.JobTypeSub {
			instance.wsClient, err = ws.Connect(ctx, ns.NodeUrlWs)
			if err != nil {
				return nil, nil, fmt.Errorf("error connect node ws %s: %v", name, err)
			}
			instance.sub, err = instance.wsClient.LogsSubscribeMentions(
				ns.Contract,
				rpc.CommitmentConfirmed,
			)
			if err != nil {
				return nil, nil, fmt.Errorf("error init log subscribe %s: %v", name, err)
			}
		}
		listeners = append(listeners, instance)
	}

	return operationRpcClient, listeners, nil
}

var cache sync.Map = sync.Map{}

type Event struct {
	EventType           string    `json:"eventType"`
	Uri                 string    `json:"uri"`
	Name                string    `json:"name"`
	MintAddress         string    `json:"mintAddress"`
	Price               float64   `json:"price"`
	Date                time.Time `json:"date"`
	Seller              string    `json:"seller"`
	AuctionHouseAddress string    `json:"auctionHouseAddress"`
	SellerReferral      string    `json:"sellerReferral"`
	Rank                int       `json:"rank"`
}

func (nl *NodeListener) HandleSignature(signature solana.Signature, data *ws.LogResult) (err error) {

	// Cache duplicated transactions
	_, loaded := nl.signatureCache.LoadOrStore(signature, "signature")
	if loaded {
		return nil
	}

	nl.Log.WithField("signature", signature).Trace("handle signature")

	go func(ctx context.Context, signature solana.Signature, nodeSettings blockchain.NodeSetting, data *ws.LogResult) {

		trxType := blockchain.TrxTypeUnknown

		// For sub job type only
		// can recognize listing in log subscription earlier, without trx data
		if nodeSettings.JobType == blockchain.JobTypeSub {
			trxType = CheckTrxType(data.Value.Logs, nodeSettings)
			if trxType == blockchain.TrxTypeUnknown {
				return
			}
		}

		trxData, err := GetTrxData(ctx, nl.rpcClient, signature, nodeSettings)
		if err != nil {
			nl.Log.Errorf("get trx data error: %s", err)
			return
		}

		// recognize listing in long time pulling later
		if nodeSettings.JobType == blockchain.JobTypeLtp {
			trxType = CheckTrxType(trxData.Meta.LogMessages, nodeSettings)
			if trxType == blockchain.TrxTypeUnknown {
				return
			}
		}

		trxDataEncoded, err := DecodeTrxData(trxData, trxType, signature)
		if err != nil {
			nl.Log.Errorf("decode tx data error: %s", err)
			return
		}

		//spew.Dump(fmt.Sprintf("signature: %s %s %g\n", signature.String(), trxDataEncoded.NftAddr, trxDataEncoded.Price))
		metadataAddress := pkg.GetSolanaMetadataAddress(solana.MustPublicKeyFromBase58(trxDataEncoded.NftAddr))

		opts := &rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
			Encoding:   "base64",
		}

		r, err := config.RpcClient.GetAccountInfoWithOpts(ctx, metadataAddress, opts)
		if err != nil {
			logrus.Errorf("error get account info: %v", err)
			return
		}

		d := r.Value.Data.GetBinary()
		dec := bin.NewBorshDecoder(d)

		var meta token_metadata.Metadata
		err = dec.Decode(&meta)
		if err != nil {
			return
		}

		metaSymbol := pkg.ClearInvisibleChars(meta.Data.Symbol)
		updateAuthority := meta.UpdateAuthority.String()
		key := updateAuthority + "__" + metaSymbol
		symbol := ""
		collection := models.Collection{}
		if ci, ok := cache.Load(key); ok {
			collection = ci.(models.Collection)
			symbol = collection.Symbol
		} else {
			collection = models.Collection{}
			if err := models.DB.NewSelect().Model(&collection).Where("meta_symbol = ?", metaSymbol).Where("update_authority = ?", updateAuthority).Scan(ctx); err == sql.ErrNoRows {
				return
			} else if err != nil {
				logrus.Errorf("error in db: %v", err)
			} else {
				symbol = collection.Symbol
				cache.Store(key, collection)
			}
		}

		item := models.Item{}
		if err := models.DB.NewSelect().Model(&item).Where("token_mint = ?", meta.Mint.String()).Scan(ctx); err != nil && err != sql.ErrNoRows {
			logrus.Errorf("error read item: %v", err)
		}

		go func(ctx context.Context, collection models.Collection, item models.Item, decoded TrxDataDecoded, meta token_metadata.Metadata, trxType string) {

			event := Event{
				Uri:                 pkg.ClearInvisibleChars(meta.Data.Uri),
				Name:                pkg.ClearInvisibleChars(meta.Data.Name),
				MintAddress:         meta.Mint.String(),
				Price:               decoded.Price,
				Date:                time.Now().UTC(),
				Seller:              decoded.Seller,
				AuctionHouseAddress: decoded.AuctionHouseAddress,
				SellerReferral:      decoded.SellerReferral,
				Rank:                item.Rank,
			}

			activityType, ok := models.TrxTypeMapping[trxType]
			if !ok {
				activityType = models.ActivityTypeSale
			}
			event.EventType = activityType

			eventData, err := json.Marshal(event)
			logrus.Infof("publish [%s]: %#v", collection.Symbol, collection)
			if err != nil {
				logrus.Errorf("error json marshal [%s]: %v", collection.Symbol, err)
			} else if _, err := config.GoCent.Publish(ctx, collection.Symbol, eventData); err != nil {
				logrus.Errorf("error publish [%s]: %v", collection.Symbol, err)
			}
		}(ctx, collection, item, trxDataEncoded, meta, trxType)

		activityType, ok := models.TrxTypeMapping[trxType]
		if !ok {
			activityType = models.ActivityTypeSale
		}
		logrus.WithField("symbol", symbol).WithField("activity_type", activityType).Info("write activity")
		activity := models.Activity{
			Name:                pkg.ClearInvisibleChars(meta.Data.Name),
			MintAddress:         meta.Mint.String(),
			Uri:                 pkg.ClearInvisibleChars(meta.Data.Uri),
			Signature:           signature.String(),
			Symbol:              symbol,
			BlockTime:           time.Now().UTC().Add(-3 * time.Second), // THis is super power YAGNI code for very very cool hackers!!!
			Price:               trxDataEncoded.Price,
			ActivityType:        activityType,
			Seller:              trxDataEncoded.Seller,
			AuctionHouseAddress: trxDataEncoded.AuctionHouseAddress,
			SellerReferral:      trxDataEncoded.SellerReferral,
		}
		if _, err := models.DB.NewInsert().Model(&activity).Exec(ctx); err != nil {
			logrus.Errorf("error in db on inser: %v", err)
			return
		}

		// Actualize listings
		if activityType == models.ActivityTypeListing || activityType == models.ActivityTypeUpdatePrice {
			listing := models.Listing{
				Signature:           signature.String(),
				Symbol:              symbol,
				BlockTime:           time.Now().UTC().Add(-3 * time.Second),
				Price:               trxDataEncoded.Price,
				Name:                pkg.ClearInvisibleChars(meta.Data.Name),
				MintAddress:         meta.Mint.String(),
				Uri:                 pkg.ClearInvisibleChars(meta.Data.Uri),
				Seller:              trxDataEncoded.Seller,
				AuctionHouseAddress: trxDataEncoded.AuctionHouseAddress,
				SellerReferral:      trxDataEncoded.SellerReferral,
				Rank:                item.Rank,
				Attributes:          item.Attributes,
			}
			query := models.DB.NewInsert().Model(&listing).On("CONFLICT(symbol, mint_address) DO UPDATE").
				Set("signature = excluded.signature").
				Set("block_time = excluded.block_time").
				Set("price = excluded.price").
				Set("name = excluded.name").
				Set("uri = excluded.uri").
				Set("seller = excluded.seller").
				Set("auction_house_address = excluded.auction_house_address").
				Set("seller_referral = excluded.seller_referral").
				Set("rank = excluded.rank").
				Set("attributes = excluded.attributes")
			if _, err := query.Exec(ctx); err != nil {
				logrus.Errorf("error insert listing: %v", err)
			}
		} else if activityType == models.ActivityTypeSale || activityType == models.ActivityTypeDelisting {
			query := models.DB.NewDelete().Model((*models.Listing)(nil)).
				Where("symbol = ?", symbol).
				Where("mint_address = ?", meta.Mint.String())
			if _, err := query.Exec(ctx); err != nil {
				logrus.Errorf("error delete listing: %v", err)
			}
		}

	}(nl.ctx, signature, nl.nodeSetting, data)

	return nil
}

func getArweave(uri string) (res ArweaveResponse, err error) {
	r, err := http.Get(uri)
	if err != nil {
		return res, err
	}
	defer r.Body.Close()
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(d, &res); err != nil {
		return res, err
	}
	return res, nil
}

func (nl *NodeListener) Listen(wg *sync.WaitGroup) {
	defer wg.Done()
	nl.Log.Info("listen")

	if nl.nodeSetting.JobType == blockchain.JobTypeSub {
		recvErr := make(chan error)
		dataCh := make(chan *ws.LogResult)
		go func(dataCh chan *ws.LogResult) {
			for {
				data, err := nl.sub.Recv()
				if err != nil {
					recvErr <- err
					return
				}
				dataCh <- data
			}
		}(dataCh)

		for {
			select {
			case err := <-recvErr:
				logrus.Errorf("node recv error: %v", err)
				nl.sub.Unsubscribe()
				nl.wsClient.Close()
				return
			case <-nl.ctx.Done():
				nl.sub.Unsubscribe()
				nl.wsClient.Close()
				return
			case data := <-dataCh:
				signature := data.Value.Signature
				err := nl.HandleSignature(signature, data)
				if err != nil {
					nl.Log.Errorf("can't handle signature: %s", signature)
					continue
				}
			}
		}

	} else if nl.nodeSetting.JobType == blockchain.JobTypeLtp {
		var lastSignature *solana.Signature

		for {
			limit := 200
			opts := &rpc.GetSignaturesForAddressOpts{
				Limit:      &limit,
				Commitment: rpc.CommitmentConfirmed,
			}
			if lastSignature != nil {
				opts.Until = *lastSignature
			}
			nl.Log.Trace("rpc - start get signatures for address")
			data, err := nl.rpcClient.GetSignaturesForAddressWithOpts(
				nl.ctx,
				nl.nodeSetting.Contract,
				opts,
			)
			if err != nil {
				nl.Log.Errorf("error get transactions: %s", err)
				continue
			}

			for j := len(data) - 1; j >= 0; j-- {
				nl.Log.WithField("i", j).Trace("range data")
				signature := data[j].Signature
				err := nl.HandleSignature(signature, nil)
				if err != nil {
					nl.Log.Errorf("can't handle signature: %s", signature)
					continue
				}
			}

			if len(data) > 0 {
				lastSignature = &data[0].Signature
			}

			time.Sleep(800 * time.Millisecond)
		}

	} else {
		nl.Log.Errorf("unknown job type: %s", nl.nodeSetting.JobType)
	}
}
