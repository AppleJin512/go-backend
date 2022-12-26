package stats

import (
	"context"
	"encoding/hex"
	"fmt"
	"moonbite/trending/internal/blockchain"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

type ArweaveResponse struct {
	Name                 string                 `json:"name"`
	Symbol               string                 `json:"symbol"`
	Description          string                 `json:"description"`
	SellerFeeBasisPoints interface{}            `json:"seller_fee_basis_points"`
	Image                string                 `json:"image"`
	ExternalURL          string                 `json:"external_url"`
	Attributes           []ArweaveResponseEvent `json:"attributes"`
	Properties           struct {
		Creators []struct {
			Address string `json:"address"`
			Share   int    `json:"share"`
		} `json:"creators"`
	} `json:"properties"`
	Collection struct {
		Name   string `json:"name"`
		Family string `json:"family"`
	} `json:"collection"`
	Signature string `json:"_signature"`
}

type ArweaveResponseEvent struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

func CheckTrxType(logs []string, settings blockchain.NodeSetting) (trxType string) {

	trxType = blockchain.TrxTypeUnknown

	if settings.Contract == blockchain.MESmartContractV2 {
		// Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K invoke [1]
		//Program log: Instruction: Deposit
		//Program 11111111111111111111111111111111 invoke [2]
		//Program 11111111111111111111111111111111 success
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K consumed 9465 of 600000 compute units
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K success
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K invoke [1]
		//Program log: Instruction: Buy
		//Program 11111111111111111111111111111111 invoke [2]
		//Program 11111111111111111111111111111111 success
		//Program log: {"price":4000000000,"buyer_expiry":0}
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K consumed 36919 of 590535 compute units
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K success
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K invoke [1]
		//Program log: Instruction: ExecuteSale
		//Program 11111111111111111111111111111111 invoke [2]
		//Program 11111111111111111111111111111111 success
		//Program 11111111111111111111111111111111 invoke [2]
		//Program 11111111111111111111111111111111 success
		//Program 11111111111111111111111111111111 invoke [2]
		//Program 11111111111111111111111111111111 success
		//Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL invoke [2]
		//Program log: Create
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [3]
		//Program log: Instruction: GetAccountDataSize
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 1622 of 445143 compute units
		//Program return: TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA pQAAAAAAAAA=
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
		//Program 11111111111111111111111111111111 invoke [3]
		//Program 11111111111111111111111111111111 success
		//Program log: Initialize the associated token account
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [3]
		//Program log: Instruction: InitializeImmutableOwner
		//Program log: Please upgrade to SPL Token 2022 for immutable owner support
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 1405 of 438653 compute units
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [3]
		//Program log: Instruction: InitializeAccount3
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 4241 of 434771 compute units
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
		//Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL consumed 24944 of 455157 compute units
		//Program ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL success
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]
		//Program log: Instruction: Transfer
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 4645 of 420019 compute units
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]
		//Program log: Instruction: CloseAccount
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 3033 of 400683 compute units
		//Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success
		//Program log: {"price":4000000000,"seller_expiry":-1,"buyer_expiry":0}
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K consumed 158887 of 553616 compute units
		//Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K success
		if len(logs) >= 40 && logs[7] == "Program log: Instruction: Buy" {
			trxType = blockchain.TrxTypeMeBuyV2
		} else if len(logs) == 8 && logs[1] == "Program log: Instruction: CancelSell" {
			trxType = blockchain.TrxTypeMeDelistingV2
		} else if len(logs) == 11 && logs[1] == "Program log: Instruction: Sell" {
			// (solana.Signature) (len=64 cap=64) 5HCWhAVnowhfLqme6h5YXzACZrcJP4wC7o78PsUaRkwinfsKYVcqMk1euwfcxsyvgbtQZ6hrxCGTjo2BBv6Zxzh4
			// ([]string) (len=11 cap=16) {
			// 	    (string) (len=62) "Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K invoke [1]",
			// 		(string) (len=30) "Program log: Instruction: Sell",
			// 		(string) (len=51) "Program 11111111111111111111111111111111 invoke [2]",
			// 		(string) (len=48) "Program 11111111111111111111111111111111 success",
			// 		(string) (len=62) "Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [2]",
			// 		(string) (len=38) "Program log: Instruction: SetAuthority",
			// 		(string) (len=89) "Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA consumed 1770 of 150523 compute units",
			// 		(string) (len=59) "Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
			// 		(string) (len=51) "Program log: {\"price\":140000000,\"seller_expiry\":-1}",
			// 		(string) (len=90) "Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K consumed 55127 of 200000 compute units",
			// 		(string) (len=59) "Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K success"
			// }
			trxType = blockchain.TrxTypeMeListingV2
		} else if len(logs) == 5 && logs[1] == "Program log: Instruction: Sell" {
			// (	[]string) (len=5 cap=8) {
			// 		(string) (len=62) "Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K invoke [1]",
			// 		(string) (len=30) "Program log: Instruction: Sell",
			// 		(string) (len=51) "Program log: {\"price\":449000000,\"seller_expiry\":-1}",
			// 		(string) (len=90) "Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K consumed 40362 of 200000 compute units",
			// 		(string) (len=59) "Program M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K success"
			// }
			trxType = blockchain.TrxTypeMeListingUpdatePriceV2
		}
	}

	//fmt.Printf("Signature type: %s\n", trxType)
	//if trxType == common.TrxTypeUnknown {
	//	spew.Dump(logs)
	//}

	return trxType
}

func GetTrxData(ctx context.Context, rpcClient *rpc.Client, signature solana.Signature, settings blockchain.NodeSetting) (trx *rpc.TransactionWithMeta, err error) {

	opts := &rpc.GetTransactionOpts{
		Encoding:   solana.EncodingBase64,
		Commitment: rpc.CommitmentConfirmed,
	}

	retryCount := 2
	for i := 0; i < retryCount; i++ {

		if i != 0 {
			time.Sleep(time.Duration(200*i) * time.Millisecond)
		}

		trx, err = rpcClient.GetConfirmedTransactionWithOpts(
			ctx,
			signature,
			opts,
		)
		if err != nil {
			//fmt.Printf("error: %v\n", err)
			continue
		}

		return trx, nil
	}

	return nil, fmt.Errorf("error get trx %s: %v\n", signature, err)
}

type TrxDataDecoded struct {
	Price               float64
	NftAddr             string
	Timestamp           time.Time
	Seller              string
	TokenATA            string
	SellerReferral      string
	AuctionHouseAddress string
}

func DecodeTrxData(trx *rpc.TransactionWithMeta, trxType string, signature solana.Signature) (trxDataDecoded TrxDataDecoded, err error) {

	if trxType == blockchain.TrxTypeMeBuyV2 {

		t, err := trx.GetTransaction()
		if err != nil {
			return TrxDataDecoded{}, err
		}
		instructionEncodedData := t.Message.Instructions[0].Data.String()
		decodedDataBytes, err := base58.Decode(instructionEncodedData)
		if err != nil {
			return trxDataDecoded, fmt.Errorf("error decode trx data, sig: %s data: %s error: %v\n", signature, instructionEncodedData, err)
		}

		instructionHex := hex.EncodeToString(decodedDataBytes)

		priceInstructionHex := instructionHex[18:]
		reversedInstructionHex := ReverseHex(priceInstructionHex)
		price, err := strconv.ParseUint(reversedInstructionHex, 16, 64)
		if err != nil {
			return trxDataDecoded, fmt.Errorf("error parse price: %v\n", err)
		}

		if len(trx.Meta.PostTokenBalances) == 0 {
			return trxDataDecoded, fmt.Errorf("error parse postTOkenBalances")
		}

		trxDataDecoded.NftAddr = trx.Meta.PostTokenBalances[0].Mint.String()
		trxDataDecoded.Price = float64(price) / float64(solana.LAMPORTS_PER_SOL)
		trxDataDecoded.Seller = t.Message.AccountKeys[0].String()
		trxDataDecoded.TokenATA = t.Message.AccountKeys[1].String()
		return trxDataDecoded, nil
	} else if trxType == blockchain.TrxTypeMeListingV2 || trxType == blockchain.TrxTypeMeListingUpdatePriceV2 {

		t, err := trx.GetTransaction()
		if err != nil {
			return TrxDataDecoded{}, err
		}
		instructionEncodedData := t.Message.Instructions[0].Data.String()
		decodedDataBytes, err := base58.Decode(instructionEncodedData)
		if err != nil {
			return trxDataDecoded, fmt.Errorf("error decode trx data, sig: %s data: %s error: %v\n", signature, instructionEncodedData, err)
		}

		instructionHex := hex.EncodeToString(decodedDataBytes)

		priceInstructionHex := instructionHex[20:36]
		reversedInstructionHex := ReverseHex(priceInstructionHex)
		price, err := strconv.ParseUint(reversedInstructionHex, 16, 64)
		if err != nil {
			return trxDataDecoded, fmt.Errorf("error parse price: %v\n", err)
		}

		if len(trx.Meta.PostTokenBalances) == 0 {
			return trxDataDecoded, fmt.Errorf("error parse postTOkenBalances")
		}

		trxDataDecoded.NftAddr = trx.Meta.PostTokenBalances[0].Mint.String()
		trxDataDecoded.Price = float64(price) / float64(solana.LAMPORTS_PER_SOL)
		trxDataDecoded.Seller = t.Message.AccountKeys[0].String()
		trxDataDecoded.TokenATA = t.Message.AccountKeys[1].String()
		trxDataDecoded.SellerReferral = t.Message.AccountKeys[t.Message.Instructions[0].Accounts[9]].String()
		trxDataDecoded.AuctionHouseAddress = t.Message.AccountKeys[t.Message.Instructions[0].Accounts[7]].String()

		return trxDataDecoded, nil
	} else if trxType == blockchain.TrxTypeMeDelistingV2 {
		trxDataDecoded.NftAddr = trx.Meta.PostTokenBalances[0].Mint.String()
		return trxDataDecoded, nil
	}

	//spew.Dump(signature)
	//spew.Dump(trx.Meta.LogMessages)
	return trxDataDecoded, fmt.Errorf("pleae add decode code for trxType='%s' in func DecodeTrxData", trxType)
}

func ReverseHex(str string) string {
	if len(str)%2 == 1 {
		str = "0" + str
	}
	result := ""
	for i := 0; i < len(str); i = i + 2 {
		result = str[i:i+2] + result
	}

	return result
}
