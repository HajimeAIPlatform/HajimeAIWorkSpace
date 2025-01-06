package solana

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/mr-tron/base58"
)

// WalletInfo 存储钱包信息的结构体
type WalletInfo struct {
	PublicKey  string
	PrivateKey string
}

// GenerateWallet 生成单个 Solana 钱包
func GenerateWallet() (*WalletInfo, error) {
	// 生成新的公私钥对
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	return &WalletInfo{
		PublicKey:  base58.Encode(pubKey),
		PrivateKey: base58.Encode(privKey),
	}, nil
}

// GenerateMultipleWallets 批量生成钱包
func GenerateMultipleWallets(count int) ([]*WalletInfo, error) {
	if count <= 0 {
		return nil, errors.New("count must be greater than zero")
	}

	wallets := make([]*WalletInfo, count)

	for i := 0; i < count; i++ {
		wallet, err := GenerateWallet()
		if err != nil {
			return nil, fmt.Errorf("failed to generate wallet %d: %w", i, err)
		}
		wallets[i] = wallet
	}

	return wallets, nil
}

// GenerateMultipleWalletsWrapper 是一个包装函数，适配 Task 的 Execute 签名
func GenerateMultipleWalletsWrapper(params ...any) {
	if len(params) == 0 {
		fmt.Println("No count parameter provided")
		return
	}
	fmt.Println("Count = ", params[0])
	fmt.Printf("Params: %v, Type of first param: %T\n", params, params[0])
	// 解包第一层
	nestedParams, ok := params[0].([]any)
	if !ok || len(nestedParams) == 0 {
		fmt.Println("Invalid or empty nested parameter")
		return
	}

	count, ok := nestedParams[0].(int)
	if !ok {
		fmt.Println("Invalid count parameter type")
		return
	}

	wallets, err := GenerateMultipleWallets(count)
	if err != nil {
		fmt.Printf("Error generating wallets: %v\n", err)
		return
	}

	// 处理生成的钱包信息
	for _, wallet := range wallets {
		fmt.Printf("Public Key: %s, Private Key: %s\n", wallet.PublicKey, wallet.PrivateKey)
	}
}
