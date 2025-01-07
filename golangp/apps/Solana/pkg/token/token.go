package token

import (
	"context"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
)

// TransferTokens transfers tokens from the sender to the receiver.
func TransferTokens(feePayer types.Account, sender types.Account, senderTokenAccount common.PublicKey, receiverTokenAccount common.PublicKey, mint common.PublicKey, amount uint64, decimals uint8) error {
	c := client.NewClient("https://broken-muddy-butterfly.solana-devnet.quiknode.pro/270ff8923ae3fcd2e905cf2dd38c6f379a317cca")

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
		return err
	}
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				token.TransferChecked(token.TransferCheckedParam{
					From:     senderTokenAccount,
					To:       receiverTokenAccount,
					Mint:     mint,
					Auth:     feePayer.PublicKey,
					Signers:  []common.PublicKey{},
					Amount:   amount,
					Decimals: decimals,
				}),
			},
		}),
		Signers: []types.Account{feePayer, feePayer},
	})
	if err != nil {
		log.Fatalf("failed to create transaction, err: %v", err)
		return err
	}

	txHash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw transaction error, err: %v\n", err)
		return err
	}

	log.Println("Transaction Hash:", txHash)
	return nil
}
