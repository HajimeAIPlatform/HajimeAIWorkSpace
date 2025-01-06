package task

import (
	"fmt"
	"hajime/golangp/apps/trading/cex/binance"
)

func TestBinanceConnectivy(args ...any) {
	fmt.Printf("TestBinanceConnectivy under runnable \n")
	binance.CheckConnectivity()
}
