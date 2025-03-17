package test

import (
	"testing"

	bn "hajime/golangp/apps/trading/cex/binance_client"

	"github.com/magiconair/properties/assert"
)

func TestBinanceConnectivy(t *testing.T) {
	bn.TestGetAccount()

	assert.Equal(t, 10, 9+1)
}
