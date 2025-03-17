package binance_client

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
)

type Client struct {
	client *binance.Client
}

func NewClient(apiKey, secretKey string) *Client {
	binance.UseTestnet = true
	var client *binance.Client
	return &Client{
		client: client,
	}
}

// ListPairs retrieves a list of all available trading pairs on Binance.
func (c *Client) ListPairs() ([]string, error) {
	exchangeInfo, err := c.client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange info: %v", err)
	}
	pairs := make([]string, len(exchangeInfo.Symbols))
	for i, symbol := range exchangeInfo.Symbols {
		pairs[i] = symbol.Symbol
	}
	return pairs, nil
}

// GetPrice fetches the current price of a specified trading pair (e.g., "BTCUSDT").
func (c *Client) GetPrice(symbol string) (string, error) {
	prices, err := c.client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to fetch price for %s: %v", symbol, err)
	}
	if len(prices) == 0 {
		return "", fmt.Errorf("no price found for symbol %s", symbol)
	}
	return prices[0].Price, nil
}

// PlaceOrder places a new order on Binance with the specified parameters.
// - symbol: Trading pair (e.g., "BTCUSDT")
// - side: Order side (binance.SideTypeBuy or binance.SideTypeSell)
// - orderType: Order type (binance.OrderTypeMarket or binance.OrderTypeLimit)
// - quantity: Amount to buy or sell (as a string for precision)
// - price: Price for limit orders (empty string for market orders)
func (c *Client) PlaceOrder(symbol string, side binance.SideType, orderType binance.OrderType, quantity string, price string) (*binance.CreateOrderResponse, error) {
	if orderType == binance.OrderTypeLimit && price == "" {
		return nil, fmt.Errorf("price is required for limit orders")
	}
	service := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(orderType).
		Quantity(quantity)
	if price != "" {
		service = service.Price(price)
	}

	order, err := service.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %v", err)
	}

	return order, nil
}

// Balance represents the balance information for an asset.
type Balance struct {
	Asset  string // The asset symbol (e.g., "BTC")
	Free   string // Available balance
	Locked string // Locked balance (e.g., in open orders)
}

// GetBalances retrieves the balances of all assets in the account.
func (c *Client) GetBalances() (map[string]Balance, error) {
	account, err := c.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account balances: %v", err)
	}
	balances := make(map[string]Balance)
	for _, b := range account.Balances {
		balances[b.Asset] = Balance{
			Asset:  b.Asset,
			Free:   b.Free,
			Locked: b.Locked,
		}
	}
	return balances, nil
}

// GetBalance retrieves the balance of a specific asset (e.g., "BTC").
func (c *Client) GetBalance(asset string) (Balance, error) {
	balances, err := c.GetBalances()
	if err != nil {
		return Balance{}, err
	}
	if balance, ok := balances[asset]; ok {
		return balance, nil
	}
	return Balance{}, fmt.Errorf("asset %s not found", asset)
}
