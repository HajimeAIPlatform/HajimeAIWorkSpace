package test

import (
	"fmt"
	"hajime/golangp/common/dify"
	"testing"
)

func TestDify(t *testing.T) {
	difyClient, err := dify.GetDifyClient()
	if err != nil {
		t.Errorf("Error getting Dify client: %v", err)
	}

	t.Logf("Dify client: %v", difyClient)

	const accessKey = "Bearer app-oMRgU9QIJKP2NBYw25EpkPyD"
	result, err := difyClient.ChatMessages("What are the specs of the iPhone 13 Pro Max?", nil, "", nil, accessKey)

	if err != nil {
		t.Errorf("Error getting chat messages: %v", err)
	}

	fmt.Printf("Chat messages: %v\n", result)
}
