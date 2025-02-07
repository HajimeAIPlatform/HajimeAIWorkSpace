package token

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
)

// TransferTokens transfers tokens from the sender to the receiver.
func TransferTokens(feePayer types.Account, sender types.Account, senderTokenAccount common.PublicKey, receiverTokenAccount common.PublicKey, mint common.PublicKey, amount uint64, decimals uint8) error {
	c := client.NewClient(GetNetworkEndpoint())

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

	LogTransaction(txHash)
	return nil
}

// TransferUSDTToToken transfers USDT to a specified token address.
//
// Parameters:
// - privateKeyHex: A hex-encoded string representing the sender's private key.
// - destTokenAddress: The public key of the destination token account.
// - mintAddress: The public key of the mint address for the token being transferred.
// - amount: The amount of tokens to transfer.
// - decimals: The number of decimal places the token uses.
//
// Returns:
// - A string representing the transaction hash if successful.
// - An error if the transaction fails.
func TransferUSDTToToken(privateKeyHex string, destTokenAddress common.PublicKey, mintAddress common.PublicKey, amount uint64, decimals uint8) (string, error) {
	client := client.NewClient(GetNetworkEndpoint())

	// Decode the private key from hex string to byte slice
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}

	// Create sender account from private key bytes
	sender, err := types.AccountFromBytes(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create sender account: %v", err)
	}

	// Get latest block hash
	recentBlockhash, err := client.GetLatestBlockhash(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get latest blockhash: %v", err)
	}

	// Create transfer instruction
	instruction := token.TransferChecked(
		token.TransferCheckedParam{
			From:     sender.PublicKey,
			To:       destTokenAddress,
			Mint:     mintAddress,
			Auth:     sender.PublicKey,
			Amount:   amount,
			Decimals: decimals,
		},
	)

	// Create transaction
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{sender},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        sender.PublicKey,
			Instructions:    []types.Instruction{instruction},
			RecentBlockhash: recentBlockhash.Blockhash,
		}),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	// Send transaction
	txHash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return txHash, nil
}

