package agent_adapter

import (
	"encoding/json"
	"fmt"
	"hajime/golangp/common/dify"
)

type DifyAdapter struct {
	AccessKey string
	Actions   []any
}

type DifyConversation struct {
	ID      string
	Adapter *DifyAdapter
	Client  *dify.DifyClient
}

func (a *DifyAdapter) CreateConversation(message string) (Conversation, error) {
	client, err := dify.GetDifyClient()
	if err != nil {
		return nil, err
	}

	for _, action := range a.Actions {
		message += "\nActions:\n"
		message += fmt.Sprintf("\n\n%s", GetFunctionDesc(action))
	}

	result, err := client.ChatMessages(message, nil, "", nil, a.AccessKey)
	if err != nil {
		return nil, err
	}

	conversation := &DifyConversation{ID: result.ConversationID, Client: client, Adapter: a}
	err = conversation.processResponse(result.Answer)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}

type Action struct {
	Name      string
	Arguments string
	Result    string
}

type Response struct {
	Actions []*Action
}

func (c *DifyConversation) processResponse(message string) error {
	var response Response
	err := json.Unmarshal([]byte(message), &response)
	if err != nil {
		return err
	}

	actionDict := map[string]any{}
	for _, action := range c.Adapter.Actions {
		actionDict[GetFunctionName(action)] = action
	}

	hasResult := false
	for _, action := range response.Actions {
		if action.Name == "" {
			continue
		}

		if actionDict[action.Name] == nil {
			fmt.Printf("Action %s not found\n", action.Name)
			continue
		}

		values := CallWithJSON(actionDict[action.Name], action.Arguments)
		result := ""
		for _, value := range values {
			result += fmt.Sprintf("%v", value)
		}
		if result != "" {
			hasResult = true
			action.Result = result
		}

	}

	if hasResult {
		jsonByte, err := json.Marshal(response)
		if err != nil {
			return err
		}
		c.SendMessage(string(jsonByte))
	}

	return nil
}

func (c *DifyConversation) SendMessage(message string) error {
	result, err := c.Client.ChatMessages(message, nil, c.ID, nil, c.Adapter.AccessKey)
	if err != nil {
		return err
	}

	err = c.processResponse(result.Answer)
	if err != nil {
		return err
	}

	return nil
}
