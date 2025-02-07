package solanaTask

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// SolanaRPC 请求结构体
type SolanaRPCRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	Id      int      `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

// SolanaRPCResponse 响应结构体
type SolanaRPCResponse struct {
	Result struct {
		TransactionCount int `json:"transactionCount"`
		Transactions     []struct {
			Slot      int    `json:"slot"`
			Signature string `json:"signature"`
			Meta      struct {
				PreBalances  []int `json:"preBalances"`
				PostBalances []int `json:"postBalances"`
			} `json:"meta"`
		} `json:"transactions"`
	} `json:"result"`
}

// Solana 地址列表
var addresses = []string{
	"BfWYgztHDqrvnf1RXGofDR49JPi7BHbkCmdGDyqSHtKe",
	"Gg8avohYTZ9G4skbqXJbtqzc99NmkTJfAhSomELWQSsh",
	"ALvmTrNzPuyJKzsHGZZBPhjCF39cXXJGhDzBwSguCGXF",
	"9npi4xTUNBwWPKaCtcWDuVVR9zsrqmWRQqLzNkeQLujG",
}

const solanaRPCURL = "https://api.mainnet-beta.solana.com"

// FetchSolanaTransactions 获取 Solana 地址的交易记录
func FetchSolanaTransactions(address string) ([]map[string]interface{}, error) {
	// 构建 Solana RPC 请求
	log.Printf("Fetching transactions for address: %s", address)
	rpcRequest := SolanaRPCRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "getConfirmedSignaturesForAddress2",
		Params:  []string{address, "finalized", "1000"}, // 获取最多 1000 签名
	}

	// 序列化请求
	reqBody, err := json.Marshal(rpcRequest)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 发送请求
	resp, err := http.Post(solanaRPCURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var solanaResp SolanaRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&solanaResp); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	log.Printf("Fetched %d transactions for address: %s", len(solanaResp.Result.Transactions), address)

	// 提取交易数据
	var transactions []map[string]interface{}
	for _, tx := range solanaResp.Result.Transactions {
		transactions = append(transactions, map[string]interface{}{
			"signature":    tx.Signature,
			"slot":         tx.Slot,
			"preBalances":  tx.Meta.PreBalances,
			"postBalances": tx.Meta.PostBalances,
		})
	}

	log.Printf("Fetched %d transactions for address: %s", len(transactions), address)
	return transactions, nil
}

// SaveTransactionsToCSV 保存交易记录到 CSV 文件
func SaveTransactionsToCSV(transactions []map[string]interface{}) error {
	log.Printf("Saving %d transactions to CSV", len(transactions))

	// 创建 CSV 文件
	filePath := "golangp/apps/hajime_center/solana_transactions.csv"
	dir := "golangp/apps/hajime_center" // 提取目录部分

	// 确保目标文件所在的目录存在
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Error creating directory: %v", err)
		return err
	}

	// 打开文件（如果文件不存在，则创建它）
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入 CSV 标题
	writer.Write([]string{"signature", "slot", "preBalances", "postBalances"})

	// 写入交易记录
	for _, tx := range transactions {
		writer.Write([]string{
			tx["signature"].(string),
			fmt.Sprintf("%d", tx["slot"].(int)),
			fmt.Sprintf("%v", tx["preBalances"]),
			fmt.Sprintf("%v", tx["postBalances"]),
		})
	}

	log.Printf("Transactions saved to CSV at: %s", filePath)
	return nil
}

// ScheduledTask 启动定时任务，定期获取 Solana 地址的交易记录并保存
func ScheduledTask() {
	// 东八区的时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("Error loading timezone: %v", err)
	}

	// 每天的固定执行时间（假设是 14:16:00）
	targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 14, 35, 0, 0, loc)

	// 如果目标时间已经过去，则设置下一次执行时间为明天的 14:16:00
	if time.Now().After(targetTime) {
		targetTime = targetTime.Add(24 * time.Hour)
	}

	// 计算距离目标时间的间隔
	durationUntilNextRun := targetTime.Sub(time.Now())

	log.Printf("Scheduled task started. First execution will be at: %s", targetTime.Format("2006-01-02 15:04:05"))

	// 定时器在目标时间触发第一次任务
	timer := time.NewTimer(durationUntilNextRun)
	defer timer.Stop()

	for {
		select {
		case <-timer.C: // 到达指定时间时执行任务
			log.Printf("Executing task at: %s", time.Now().Format("2006-01-02 15:04:05"))

			// 执行任务
			for _, address := range addresses {
				transactions, err := FetchSolanaTransactions(address)
				if err != nil {
					log.Printf("Error fetching transactions for address %s: %v", address, err)
					continue
				}

				if err := SaveTransactionsToCSV(transactions); err != nil {
					log.Printf("Error saving transactions to CSV for address %s: %v", address, err)
				}
			}

			// 重置定时器，计算下一天的执行时间
			targetTime = targetTime.Add(24 * time.Hour) // 每天执行一次
			durationUntilNextRun = targetTime.Sub(time.Now())
			timer.Reset(durationUntilNextRun)
		}
	}
}
