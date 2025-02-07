package task

import (
	"fmt"
	"hajime/golangp/apps/trading/cex/binance"
	"sync"
)

func TestBinanceConnectivy() {
	fmt.Printf("TestBinanceConnectivy under runnable \n")
	binance.CheckConnectivity()
}

func CheckBinanceMarketData(symbol string, wg *sync.WaitGroup) {
	fmt.Printf("CheckBinanceMarketData under runnable for symbol: %s\n", symbol)
	binance.CheckMarketPrice(symbol)
	wg.Done()
}
