package pkg

import "github.com/gagliardetto/solana-go"

var TokenMetadataProgramID = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

func GetSolanaMetadataAddress(mint solana.PublicKey) solana.PublicKey {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			TokenMetadataProgramID.Bytes(),
			mint.Bytes(),
		},
		TokenMetadataProgramID,
	)
	if err != nil {
		panic(err)
	}
	return addr
}
