package token

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/blocto/solana-go-sdk/client"
)

var NetworkConfig = map[string]string{
	"mainnet":            "https://api.mainnet-beta.solana.com",
	"quick_node_mainnet": "https://broken-muddy-butterfly.solana-devnet.quiknode.pro/270ff8923ae3fcd2e905cf2dd38c6f379a317cca/",
	"testnet":            "https://api.testnet.solana.com",
	"devnet":             "https://api.devnet.solana.com",
	"quick_node_devnet":  "https://attentive-ancient-sun.solana-devnet.quiknode.pro/10516a7a532763abe977fe3be9b3fc99cd6a5453/",
	"localhost":          "http://localhost:8899",
}

// 获取环境变量中的网络地址
func GetNetworkEndpoint() string {
	network := os.Getenv("SOLANA_NETWORK")
	if network == "" {
		network = NetworkConfig["quick_node_devnet"]
	}
	return network
}

// 保存交易哈希并打印链接
func LogTransaction(txHash string) {
	log.Printf("Transaction hash: %s", txHash)
	fmt.Printf("View transaction: https://explorer.solana.com/tx/%s\n", txHash)
}

// 获取Token余额
func GetTokenBalanceSpl(publicKey string) (uint64, error) {
	client := client.NewClient(GetNetworkEndpoint())

	balance, err := client.GetTokenAccountBalance(context.Background(), publicKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get token balance: %v", err)
	}

	return balance.Amount, nil
}
