package task

import (
	"fmt"
	"hajime/golangp/apps/trading/cex/binance_client"
	"sync"
)

func TestBinanceConnectivy() {
	fmt.Printf("TestBinanceConnectivy under runnable \n")
	binance_client.CheckConnectivity()
}

func CheckBinanceMarketData(symbol string, wg *sync.WaitGroup) {
	fmt.Printf("CheckBinanceMarketData under runnable for symbol: %s\n", symbol)
	binance_client.CheckMarketPrice(symbol)
	wg.Done()
}
