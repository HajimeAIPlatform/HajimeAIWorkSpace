package test

import (
	"fmt"
	"hajime/golangp/apps/AgentTown/agent_adapter"
	"hajime/golangp/apps/trading/dex/solana/raydium"
	"testing"
)

type RaydiumClient struct {
	privateKey string
}

func (r *RaydiumClient) MakeTrade(p struct {
	TokenIn  string `desc:"Token symbol to swap from"`
	TokenOut string `desc:"Token symbol to swap to"`
	AmountIn int64  `desc:"Amount of token to swap"`
}) string {
	// return fmt.Sprintf("Trade %d %s for %s", p.AmountIn, p.TokenIn, p.TokenOut)
	res, err := raydium.CallSwap(
		p.TokenIn,
		p.TokenOut,
		r.privateKey,
		p.AmountIn,
		0)
	if err != nil {
		return fmt.Sprintf("Error making trade: %v", err)
	}
	return res
}

type Store struct {
	USDCBalance  int
	TRUMPBalance int
	Memory       string
}

func (s *Store) UpdateMemory(p struct {
	Memory string `desc:"Memory to append"`
}) {
	s.Memory += p.Memory
	s.Memory += "\n"
}

func (s *Store) UpdateUSDCBalance(p struct {
	Balance int `desc:"Balance to update"`
}) {
	s.USDCBalance = p.Balance
}

func (s *Store) UpdateTRUMPBalance(p struct {
	Balance int `desc:"Balance to update"`
}) {
	s.TRUMPBalance = p.Balance
}

func EndConversation(p struct{}) string {
	return agent_adapter.END_CONVERSATION
}

func TestTradingAgent(t *testing.T) {
	var raydiumClient = &RaydiumClient{
		privateKey: "2dGuRnMehWAr6pXg3A76ivkFZQDYN3HivC9pqtteN69SZmfaY8KBkyqATgtJhgDJeUEREzbrkEzg4cfboesAbKjz",
	}

	var store = &Store{
		USDCBalance: 10000,
	}

	var adapter agent_adapter.Adapter = &agent_adapter.DifyAdapter{
		AccessKey: "Bearer app-RaG6pT0giNLBm3DmWXlPRWPL",
		Actions: []any{
			raydiumClient.MakeTrade,
			store.UpdateUSDCBalance,
			store.UpdateTRUMPBalance,
			store.UpdateMemory,
			EndConversation,
		},
	}

	// Demo begin
	dates := []string{"2025-10-01", "2025-10-02", "2025-10-03", "2025-10-04", "2025-10-05"}
	prices := []int{20, 15, 25, 30, 20}

	for i, price := range prices {
		date := dates[i]
		_, err := adapter.CreateConversation(fmt.Sprintf(`
		System:
		You are a professional trading agent, 
		Use the action provided to trade USDC/TRUMP pair.
		Actions will be executed in parallel.
		You can also store the trade information in memory.
		Trade once perday.
		Today's date is %s.
		Today's price for TRUMP is %d USDC.
		
		Memory: 
		%s

		Balance:
		USDC: %d
		TRUMP: %d
		`, date, price, store.Memory, store.USDCBalance, store.TRUMPBalance))

		if err != nil {
			t.Errorf("Error starting conversation: %v", err)
		}
	}
}

func TestSwap(t *testing.T) {
	var raydiumClient = &RaydiumClient{
		privateKey: "2dGuRnMehWAr6pXg3A76ivkFZQDYN3HivC9pqtteN69SZmfaY8KBkyqATgtJhgDJeUEREzbrkEzg4cfboesAbKjz",
	}

	res := raydiumClient.MakeTrade(struct {
		TokenIn  string `desc:"Token symbol to swap from"`
		TokenOut string `desc:"Token symbol to swap to"`
		AmountIn int64  `desc:"Amount of token to swap"`
	}{
		TokenIn:  "TRUMP",
		TokenOut: "USDC",
		AmountIn: 236,
	})

	fmt.Println(res)
}
