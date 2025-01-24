package main

import (
	"fmt"
	"log"

	"hajime/golangp/apps/raydium"
)

func main() {
	tokenIn := "SOL"
	tokenOut := "USDT"
	privateKey := ""
	amountIn := int64(100000)
	microLamports := int64(300000)

	txId, err := raydium.CallSwap(tokenIn, tokenOut, privateKey, amountIn, microLamports)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Swap successful. TxId: %s\n", txId)
}
