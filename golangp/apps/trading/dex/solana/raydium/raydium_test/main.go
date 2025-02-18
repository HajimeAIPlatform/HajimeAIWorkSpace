package main

import (
	"fmt"
	"log"
	"math"

	"hajime/golangp/apps/trading/dex/solana/raydium"
)

func TestSwap() {
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

func TestCreateToken() {
	privateKey := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	tokenName := "HAJIME_F"
	tokenSymbol := "HAJIME_F"
	description := "HAJIME_M.\n\nTELEGRAM:  \nTWITTER:  \n WEBSITE: "
	uri := "https://devixyz.github.io/telegram/hajime.json"
	tokenSupply := int64(1_500_000)
	tokenDecimals := int64(6)
	Data, err := raydium.CallCreateToken(privateKey, tokenName, tokenSymbol, description, uri, tokenSupply, tokenDecimals)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("create token successful. Data: %s\n", Data)
}

func TestCreateMarket() {
	privateKey := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	mintAAddress := "3uAv9qSsUdz2RFkVx99Fe81dqpDUbxQRzN1kNwdykTuf"
	mintADecimals := 6
	mintBAddress := "So11111111111111111111111111111111111111112"
	mintBDecimals := 9

	Data, err := raydium.CallCreateMarket(privateKey, mintAAddress, mintADecimals, mintBAddress, mintBDecimals)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("create market successful. Data: %s\n", Data)
}

func TestCreatePool() {
	privateKey := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	mintAAddress := "3uAv9qSsUdz2RFkVx99Fe81dqpDUbxQRzN1kNwdykTuf"
	mintADecimals := 6
	mintAInitialAmount := int64(1_000 * math.Pow(10, 6))
	mintBAddress := "So11111111111111111111111111111111111111112"
	mintBDecimals := 9
	mintBInitialAmount := int64(1 * math.Pow(10, 9))
	marketId := "DRSbrtzZPoAwwJ36kEfRfQka5jaBmfcm46NmBt5ASnSu"

	Data, err := raydium.CallCreatePool(privateKey, mintAAddress, mintADecimals, mintAInitialAmount, mintBAddress, mintBDecimals, mintBInitialAmount, marketId)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("create market successful. Data: %s\n", Data)
}

func main() {
	// TestSwap()
	// TestCreateToken()
	// TestCreateMarket()
	TestCreatePool()
}
