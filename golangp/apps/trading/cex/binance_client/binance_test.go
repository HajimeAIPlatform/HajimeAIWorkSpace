package binance_client_test

import (
	"context"
	"fmt"
	"hajime/golangp/apps/trading/cex/binance_client"
	"log"
	"testing"

	"github.com/adshao/go-binance/v2"
)

func TestBClient(t *testing.T) {
	binance.UseTestnet = true
	client := binance.NewClient("Izl2SNYJB3hGqAp4jFk41InSLt2KHMRDbG3BkeaxeOOgcMUFNSjO09Q2yRoTqomy", "ydakwvYpMcbyf6dNXwGE4Wtxxz1OpB1IGKiUF1i05DLLY1wb5cN9EJZa1PXBhErK")
	// Create a context
	ctx := context.Background()

	// Fetch account balance
	account, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching balance: %v", err)
	}

	// Print the balances
	fmt.Println("Account Balances:")
	for _, balance := range account.Balances {
		// Only display assets with a non-zero balance
		fmt.Printf("Asset: %s, Balance: %s, Available: %s\n",
			balance.Asset,
			balance.Free,
			balance.Locked)
	}

	price, err := client.NewListPricesService().Symbol("BTCUSDT").Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching price: %v", err)
	}
	fmt.Println("BTCUSDT Price:")
	for _, p := range price {
		fmt.Printf("Symbol: %s, Price: %s\n", p.Symbol, p.Price)
	}
}

func TestClient(t *testing.T) {
	client := binance_client.NewClient("Izl2SNYJB3hGqAp4jFk41InSLt2KHMRDbG3BkeaxeOOgcMUFNSjO09Q2yRoTqomy", "ydakwvYpMcbyf6dNXwGE4Wtxxz1OpB1IGKiUF1i05DLLY1wb5cN9EJZa1PXBhErK")
	balances, err := client.GetBalances()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Balances:", balances)
	// List trading pairs
	pairs, err := client.ListPairs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Trading pairs:", pairs[:5]) // Print first 5 pairs

	// Check price
	price, err := client.GetPrice("BTCUSDT")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("BTCUSDT Price:", price)

	// Place a market buy order
	order, err := client.PlaceOrder("BTCUSDT", binance.SideTypeBuy, binance.OrderTypeMarket, "0.001", "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Order ID:", order.OrderID)

	// Check balance
	balance, err := client.GetBalance("BTC")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("BTC Balance - Free: %s, Locked: %s\n", balance.Free, balance.Locked)
}
