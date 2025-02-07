package test

import (
	"hajime/golangp/apps/Solana/pkg/token"
	"testing"

	"github.com/blocto/solana-go-sdk/common"
)

func TestTransferUSDTToToken(t *testing.T) {
	// Replace with actual test values
	privateKeyHex := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	destTokenAddress := common.PublicKeyFromString("destination_token_address")
	mintAddress := common.PublicKeyFromString("2wkoqByNMi3dSDJbeH8WtWhXcogFLZvEWzzGyJJbCnCc")
	amount := uint64(1000000) // Example amount
	decimals := uint8(6)      // Example decimals

	txHash, err := token.TransferUSDTToToken(privateKeyHex, destTokenAddress, mintAddress, amount, decimals)
	if err != nil {
		t.Fatalf("TransferUSDTToToken failed: %v", err)
	}

	if txHash == "" {
		t.Fatalf("Expected a transaction hash, got an empty string")
	}

	t.Logf("Transaction successful, hash: %s", txHash)
}
