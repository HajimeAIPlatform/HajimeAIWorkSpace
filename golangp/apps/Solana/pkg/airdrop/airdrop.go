package airdrop

import (
	"context"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/types"
)

// CheckAndAirdrop checks the balance of an account and requests an airdrop if necessary.
func CheckAndAirdrop(account types.Account, accountName string) {
	c := client.NewClient("https://broken-muddy-butterfly.solana-devnet.quiknode.pro/270ff8923ae3fcd2e905cf2dd38c6f379a317cca")

	balance, err := c.GetBalance(context.Background(), account.PublicKey.ToBase58())
	if err != nil {
		log.Fatalf("Error checking balance for %s: %v", accountName, err)
	}
	fmt.Printf("Balance for %s: %d lamports\n", accountName, balance)

	if balance == 0 {
		txHash, err := c.RequestAirdrop(context.Background(), account.PublicKey.ToBase58(), 1000000000)
		if err != nil {
			log.Fatalf("Error requesting airdrop for %s: %v", accountName, err)
		}
		fmt.Printf("Airdrop successful for %s, txHash: %s\n", accountName, txHash)
	}
}
