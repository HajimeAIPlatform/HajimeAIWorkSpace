package main

import (
	"fmt"
	"log"

	"hajime/golangp/apps/raydium"
)

func main() {
	tokenIn := "SOL"
	tokenOut := "USDT"
	privateKey := "3U4fCwNpeH4WXU4MTWDoBjS2ps9pm5y74syuZKq9PHMYEh7jMxcTwvhvdmz7Ee4cj5ANL18Q5ceSU4gQWxPkFuQY"
	amountIn := int64(10000)
	microLamports := int64(300000)

	txId, err := raydium.CallSwap(tokenIn, tokenOut, privateKey, amountIn, microLamports)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Swap successful. TxId: %s\n", txId)
}
