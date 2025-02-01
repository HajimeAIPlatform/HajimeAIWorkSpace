package test

import (
	"fmt"
	"hajime/golangp/apps/AgentTown/agent_adapter"
	"testing"
)

type RaydiumClient struct {
	privateKey string
}

// Mock trade api
func (r *RaydiumClient) MakeTrade(p struct {
	TokenIn       string `desc:"Token to swap from"`
	TokenOut      string `desc:"Token to swap to"`
	AmountIn      int64  `desc:"Amount of token to swap"`
	MicroLamports int64  `desc:"Amount of token to swap in micro lamports"`
}) string {
	fmt.Printf("Executing solana swap from %s to %s\n", p.TokenIn, p.TokenOut, p.AmountIn, p.MicroLamports)
	return "The swap was successful"
}

func CheckNews(p struct {
	Keyword string `desc:"Keyword to search for"`
}) string {
	fmt.Printf("Checking news for keyword: %s\n", p.Keyword)
	return "World War 3 is about to start"
}

func TestTradingAgent(t *testing.T) {
	var raydiumClient = &RaydiumClient{
		privateKey: "YOUR_PRIVATE",
	}

	var adapter agent_adapter.Adapter = &agent_adapter.DifyAdapter{
		AccessKey: "Bearer app-RaG6pT0giNLBm3DmWXlPRWPL",
		Actions: []any{
			raydiumClient.MakeTrade,
			CheckNews,
		},
	}

	_, err := adapter.CreateConversation(`
	You are a trading agent, you take the necessary actions provided to solve problem.
	You can conclude the task if you finish the task or you require more information.

	User Input:
	Hi, I heard about today $Trump is sky rocketing, 
	I want to become rich please trade for me.
	`)

	if err != nil {
		t.Errorf("Error starting conversation: %v", err)
	}
}
