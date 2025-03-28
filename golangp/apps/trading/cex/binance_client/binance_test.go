package binance_client_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
)

func strSliceToMap(strings []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range strings {
		result[s] = true
	}
	return result
}

func TestListAsset(t *testing.T) {
	binance.UseTestnet = true
	client := binance.NewClient("Izl2SNYJB3hGqAp4jFk41InSLt2KHMRDbG3BkeaxeOOgcMUFNSjO09Q2yRoTqomy", "ydakwvYpMcbyf6dNXwGE4Wtxxz1OpB1IGKiUF1i05DLLY1wb5cN9EJZa1PXBhErK")
	// Create a context
	ctx := context.Background()

	interestedAssets := []string{"BTC", "ETH", "SOL", "USDT"}
	interestedAssetMap := strSliceToMap(interestedAssets)

	account, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching balance: %v", err)
	}

	// Print the balances for interested assets only
	fmt.Println("Account Balances (Interested Assets Only):")
	for _, balance := range account.Balances {
		// Check if the asset is in our interested list and has a non-zero balance
		if interestedAssetMap[balance.Asset] {
			// Convert string balances to float for checking non-zero
			free, _ := strconv.ParseFloat(balance.Free, 64)
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			if free > 0 || locked > 0 {
				fmt.Printf("Asset: %s, Balance: %s, Locked: %s\n",
					balance.Asset,
					balance.Free,
					balance.Locked)
			}
		}
	}
}

func TestCheckPrice(t *testing.T) {
	binance.UseTestnet = true
	client := binance.NewClient("Izl2SNYJB3hGqAp4jFk41InSLt2KHMRDbG3BkeaxeOOgcMUFNSjO09Q2yRoTqomy", "ydakwvYpMcbyf6dNXwGE4Wtxxz1OpB1IGKiUF1i05DLLY1wb5cN9EJZa1PXBhErK")
	// Create a context
	ctx := context.Background()
	interestedPairs := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}

	// Calculate time range for last 7 days
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7) // 7 days ago

	for _, pair := range interestedPairs {
		// Get price data for last 7 days for each pair
		klines, err := client.NewKlinesService().
			Symbol(pair).
			Interval("1d").                     // Daily candles
			StartTime(startTime.Unix() * 1000). // Convert to milliseconds
			EndTime(endTime.Unix() * 1000).
			Do(ctx)
		if err != nil {
			log.Fatalf("Error fetching klines for %s: %v", pair, err)
		}

		fmt.Printf("%s Prices (Last 7 Days):\n", pair)
		for _, k := range klines {
			openPrice, _ := strconv.ParseFloat(k.Open, 64)
			highPrice, _ := strconv.ParseFloat(k.High, 64)
			lowPrice, _ := strconv.ParseFloat(k.Low, 64)
			closePrice, _ := strconv.ParseFloat(k.Close, 64)
			timestamp := time.Unix(k.CloseTime/1000, 0) // Convert milliseconds to seconds

			fmt.Printf("Date: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f\n",
				timestamp.Format("2006-01-02"),
				openPrice,
				highPrice,
				lowPrice,
				closePrice)
		}

		fmt.Println() // Add a blank line between pairs
	}
}

func TestBuyAndSell(t *testing.T) {
	binance.UseTestnet = true
	client := binance.NewClient("Izl2SNYJB3hGqAp4jFk41InSLt2KHMRDbG3BkeaxeOOgcMUFNSjO09Q2yRoTqomy", "ydakwvYpMcbyf6dNXwGE4Wtxxz1OpB1IGKiUF1i05DLLY1wb5cN9EJZa1PXBhErK")
	ctx := context.Background()

	// Define the trading pair and quantity
	symbol := "BTCUSDT"
	quantity := "0.001" // Small quantity for testing (0.001 BTC)

	// Step 1: Check initial balance
	account, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching initial balance: %v", err)
	}
	initialUSDT := getBalance(account, "USDT")
	initialBTC := getBalance(account, "BTC")
	fmt.Printf("Initial Balances - USDT: %s, BTC: %s\n", initialUSDT.Free, initialBTC.Free)

	// Step 2: Place a market buy order
	buyOrder, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeMarket).
		Quantity(quantity).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error placing buy order: %v", err)
	}
	fmt.Printf("Buy Order Executed - Order ID: %d, Filled: %s\n", buyOrder.OrderID, buyOrder.ExecutedQuantity)

	// Step 3: Wait briefly and check balance after buy
	time.Sleep(2 * time.Second) // Wait for order to process
	account, err = client.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching balance after buy: %v", err)
	}
	afterBuyUSDT := getBalance(account, "USDT")
	afterBuyBTC := getBalance(account, "BTC")
	fmt.Printf("After Buy Balances - USDT: %s, BTC: %s\n", afterBuyUSDT.Free, afterBuyBTC.Free)

	// Step 4: Place a market sell order
	sellOrder, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeSell).
		Type(binance.OrderTypeMarket).
		Quantity(quantity).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error placing sell order: %v", err)
	}
	fmt.Printf("Sell Order Executed - Order ID: %d, Filled: %s\n", sellOrder.OrderID, sellOrder.ExecutedQuantity)

	// Step 5: Check final balance
	time.Sleep(2 * time.Second) // Wait for order to process
	account, err = client.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching final balance: %v", err)
	}
	finalUSDT := getBalance(account, "USDT")
	finalBTC := getBalance(account, "BTC")
	fmt.Printf("Final Balances - USDT: %s, BTC: %s\n", finalUSDT.Free, finalBTC.Free)
}

// Helper function to get balance for a specific asset
func getBalance(account *binance.Account, asset string) *binance.Balance {
	for _, balance := range account.Balances {
		if balance.Asset == asset {
			return &balance
		}
	}
	return &binance.Balance{Asset: asset, Free: "0", Locked: "0"}
}
