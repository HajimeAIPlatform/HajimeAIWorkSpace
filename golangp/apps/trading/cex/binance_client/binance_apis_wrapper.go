package binance_client

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
)

func TestGetAccount() {
	// Create a new Binance futures client
	futures.UseTestnet = true
	ApiKey := ""
	SecretKey := ""
	fc := futures.NewClient(ApiKey, SecretKey)
	fmt.Printf("Binance futures client created\n")
	res, err := fc.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func CheckMarketPrice(symbol string) {
	// Create a new Binance futures client
	futures.UseTestnet = true
	ApiKey := ""
	SecretKey := ""
	fc := futures.NewClient(ApiKey, SecretKey)
	fmt.Printf("Binance futures client created\n")

	// Get the market price for the given symbol
	price, err := fc.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, p := range price {
		fmt.Printf("Symbol: %s, Price: %s\n", p.Symbol, p.Price)
	}
}

// ...existing code...
func PlaceMarketOrder(symbol, side, quantity string) error {
	futures.UseTestnet = true
	ApiKey := ""
	SecretKey := ""
	fc := futures.NewClient(ApiKey, SecretKey)

	_, err := fc.NewCreateOrderService().
		Symbol(symbol).
		Side(futures.SideType(side)).
		Type(futures.OrderTypeMarket).
		Quantity(quantity).
		Do(context.Background())
	return err
}

// ...existing code...
