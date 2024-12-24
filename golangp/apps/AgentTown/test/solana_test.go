package test

import (
	"testing"

	"hajime/golangp/apps/AgentTown/solana"
)

func TestGenerateWallet(t *testing.T) {
	wallet, err := solana.GenerateWallet()
	if err != nil {
		t.Fatalf("Failed to generate wallet: %v", err)
	}

	if wallet.PublicKey == "" || wallet.PrivateKey == "" {
		t.Errorf("Generated wallet has empty public or private key")
	}
}

func TestGenerateMultipleWallets(t *testing.T) {
	count := 5
	wallets, err := solana.GenerateMultipleWallets(count)
	if err != nil {
		t.Fatalf("Failed to generate multiple wallets: %v", err)
	}

	if len(wallets) != count {
		t.Errorf("Expected %d wallets, but got %d", count, len(wallets))
	}

	for _, wallet := range wallets {
		if wallet.PublicKey == "" || wallet.PrivateKey == "" {
			t.Errorf("Generated wallet has empty public or private key")
		}
	}
}

func TestNewConnection(t *testing.T) {
	conn := solana.NewConnection("")
	if conn == nil {
		t.Errorf("Failed to create new connection")
	} else if conn.Client == nil {
		t.Errorf("Connection client is nil")
	}
}

func TestGetWalletInfo(t *testing.T) {
	conn := solana.NewConnection("")

	// 使用一个已知的公钥进行测试，例如 Solana 的系统程序公钥
	publicKey := "11111111111111111111111111111111"
	err := conn.GetWalletInfo(publicKey)
	if err != nil {
		t.Errorf("Failed to get wallet info: %v", err)
	}
}


