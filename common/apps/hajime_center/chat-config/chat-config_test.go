package chat_config

import (
	"testing"
)

func TestReadChatConfig(t *testing.T) {
	CLI.Config = "chat-config.example.json"

	OpenaiConfig := LoadChatConfig()
	for _, element := range OpenaiConfig.AdminEmail {
		println(element)
	}
}
