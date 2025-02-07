/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-02-07 16:19:15
 */
package main

import (
	"hajime/golangp/apps/Solana/pkg/account"
	"hajime/golangp/apps/Solana/pkg/airdrop"
	"hajime/golangp/apps/Solana/pkg/token"
	"log"
	"path/filepath"

	"github.com/blocto/solana-go-sdk/common"
)

const (
	accountKeyPath1 = "golangp/apps/Solana/assets/wallet_7dEc3i8Niz.json"
	accountKeyPath2 = "golangp/apps/Solana/assets/wallet_3qQEWctNXM.json"
	transferAmount  = 100000000
)

var (
	tokenMintPubkey            = common.PublicKeyFromString("3bsQNidmWGYJZ4W8d1AtDPKomZC3RCgrfssm7f14dAeH")
	senderTokenAccountPubkey   = common.PublicKeyFromString("rXCWwt3Nx9coEubo2Xxh4bPyU6hi3pCgFzHzAsybaKH")
	receiverTokenAccountPubkey = common.PublicKeyFromString("B5fyaWUMLd2GAxEmQSEqcZ874WwEY6tfd5X6ZRQzzhC6")
)

func main() {
	dirPath := filepath.Join("/home/lio/hajime", "golangp", "apps", "Solana", "assets")

	publicKeys, err := account.CreateAccounts(dirPath, 5)
	log.Println(publicKeys)
	if err != nil {
		log.Fatalf("failed to create accounts: %v", err)
	}

	feePayerAccount, err := account.LoadAccountFromFile(accountKeyPath1)
	if err != nil {
		log.Fatalf("failed to load feePayer account: %v", err)
	}
	senderAccount, err := account.LoadAccountFromFile(accountKeyPath2)
	if err != nil {
		log.Fatalf("failed to load sender account: %v", err)
	}

	airdrop.CheckAndAirdrop(feePayerAccount, "Fee Payer")
	airdrop.CheckAndAirdrop(senderAccount, "Sender")

	err = token.TransferTokens(
		feePayerAccount,
		senderAccount,
		senderTokenAccountPubkey,
		receiverTokenAccountPubkey,
		tokenMintPubkey,
		transferAmount,
		9,
	)
	if err != nil {
		log.Fatalf("Error transferring tokens: %v", err)
	}
}
