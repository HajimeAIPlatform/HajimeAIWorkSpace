package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenClaimController struct {
	CSVFilePath     string
	USDCPrice       float64
	LastPriceUpdate time.Time
}

// Solana地址查询的结构体
type SolanaData struct {
	USDCAmount       float64    `json:"usdc_amount"`
	PriceUSD         float64    `json:"price_usd"`
	TransactionCount int        `json:"transaction_count"`
	Activities       []Activity `json:"activities"`
}

type Activity struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Total         float64 `json:"total"`
	ProductName   string  `json:"product_name"` // Added ProductName
	Date          string  `json:"date"`
}

// Product information
var productMap = map[string]struct {
	Name  string
	Price float64
}{
	"8QAG8Zjv8xYdog8QRmiC3WcK9z3bC2fKREqVojuksvPT": {
		Name:  "Hajime Node Bot",
		Price: 1499,
	},
	"A1P2DLa8bYuu2ouCBEGGStb7AN6HtLH6zz2SpSEsDqC": {
		Name:  "Hajime Node C",
		Price: 20000,
	},
	"DprENQRtuDmBBbnxfxWPG5NsHDTiihvvc29X9TgD6pFJ": {
		Name:  "Hajime Node B",
		Price: 50000,
	},
	"B1u1cs45KFDqYBcoCHJnXgQM368qPWVaTi2XNLRE5MmK": {
		Name:  "Hajime Node A",
		Price: 100000,
	},
	"HajimeTestDqYBcoCHJnXgQM368qPWVaTi2XNLRE5MmM": {
		Name:  "Hajime Node Test",
		Price: 100000,
	},
}

// NewTokenClaimController creates a new TokenClaimController
func NewTokenClaimController(csvFilePath string) TokenClaimController {
	return TokenClaimController{
		CSVFilePath: csvFilePath,
	}
}

// GetSolanaAddressInfo retrieves the Solana address information
func (tc *TokenClaimController) GetSolanaAddressInfo(ctx *gin.Context) {
	address := ctx.Param("address")

	// Get Solana data from CSV
	solanaData, err := tc.getSolanaDataFromCSV(address)
	if err != nil {
		log.Printf("Error getting Solana data from CSV: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Solana data from CSV"})
		return
	}

	// Get USDC price
	if time.Since(tc.LastPriceUpdate).Hours() > 24 {
		priceUSD, err := tc.getUSDCPrice()
		if err != nil {
			log.Printf("Error getting USDC price: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get USDC price"})
			return
		}
		tc.USDCPrice = priceUSD
		tc.LastPriceUpdate = time.Now()
	}

	// Calculate token amount (USDC)
	tokenAmount := (solanaData.USDCAmount / 30_000_000.0) * 1e9
	tokenAmount = truncateToSixDecimals(tokenAmount)

	// Calculate total USD value
	totalUSDValue := solanaData.USDCAmount * tc.USDCPrice
	totalUSDValue = truncateToSixDecimals(totalUSDValue)

	// Return the data
	ctx.JSON(http.StatusOK, gin.H{
		"token_amount": tokenAmount,
		"price_usd":    totalUSDValue,
		"activities":   solanaData.Activities,
	})
}

// golangp/apps/hajime_center/solana__transactions.csv
// getSolanaDataFromCSV reads the CSV file and retrieves the data for a specific address
func (tc *TokenClaimController) getSolanaDataFromCSV(address string) (*SolanaData, error) {
	file, err := os.Open(tc.CSVFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// Read CSV records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	var totalAmount float64
	var activities []Activity

	// Process CSV records
	for _, record := range records {
		fromAddress := record[4]
		toAddress := record[5]
		amount, err := parseFloat(record[8])
		if err != nil {
			log.Printf("Invalid amount format in record %v: %v", record, err)
			continue
		}
		transactionID := record[1]
		blockTime := record[2] // block_time field needs conversion to date

		// Check if the address matches and process the transaction
		if fromAddress == address {
			amount = amount / 1e6 // Adjust precision
			log.Printf("address: %v", address)

			// Find product details based on the address
			product, found := productMap[toAddress]
			log.Printf("address: %v", product)

			if !found {
				continue // Skip if no matching product found
			}

			// Calculate the amount of products purchased (e.g., USDC amount / product price)
			quantity := amount / product.Price

			// Convert block time to date
			transactionDate := time.Unix(int64(parseInt(blockTime)), 0).Format("2006-01-02 15:04:05")

			// Record the transaction activity
			activities = append(activities, Activity{
				TransactionID: transactionID,
				Amount:        quantity,
				Total:         amount,
				ProductName:   product.Name,
				Date:          transactionDate,
			})

			// Sum up the total USDC amount
			totalAmount += amount
		}
	}

	// If no activities found, return an error
	if len(activities) == 0 {
		return nil, fmt.Errorf("address not found in CSV")
	}

	// Truncate total amount before returning
	totalAmount = truncateToSixDecimals(totalAmount)
	return &SolanaData{
		USDCAmount:       totalAmount,
		TransactionCount: len(activities),
		Activities:       activities,
	}, nil
}

// truncateToSixDecimals truncates a float to 6 decimal places
func truncateToSixDecimals(value float64) float64 {
	return math.Round(value*1e6) / 1e6
}

// getUSDCPrice retrieves the current USDC price (mocked for now)
func (tc *TokenClaimController) getUSDCPrice() (float64, error) {
	// resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=usd-coin&vs_currencies=usd")
	// if err != nil {
	// 	return 0.96, fmt.Errorf("failed to get USDC price from API: %w", err)
	// }
	// defer resp.Body.Close()

	// // 解析响应中的 JSON 数据
	// var result map[string]map[string]float64
	// decoder := json.NewDecoder(resp.Body)
	// err = decoder.Decode(&result)
	// if err != nil {
	// 	return 0.96, fmt.Errorf("failed to parse USDC price JSON response: %w", err)
	// }

	// // 获取 USDC 价格
	// price, exists := result["usd-coin"]["usd"]
	// if !exists {
	// 	return 0.96, fmt.Errorf("USDC price not found in response")
	// }
	price := 1.0
	return price, nil
}

// parseFloat converts a string to float64
func parseFloat(value string) (float64, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
		return 0.0, err
	}
	return result, nil
}

// parseInt converts a string to int64
func parseInt(value string) int64 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("Error parsing int: %v", err)
	}
	return result
}
