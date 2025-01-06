package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
)

func TestGetAccount() {
	// Create a new Binance futures client
	futures.UseTestnet = true
	ApiKey := "a2a7e65b0ccf7d4355074bcb1d1e29456d9fd518abbdaac308da4191cdfe4038"
	SecretKey := "7eed1bea32d0cb8991ab9c939f0c40c11a4bd7046fc68f750a9073941148b3c9"
	fc := futures.NewClient(ApiKey, SecretKey)
	fmt.Printf("Binance futures client created\n")
	res, err := fc.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
