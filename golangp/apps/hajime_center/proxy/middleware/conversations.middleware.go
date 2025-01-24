package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"net/http"

	"gorm.io/gorm"
)

type ConversationResponse struct {
	Limit   int                `json:"limit"`
	HasMore bool               `json:"has_more"`
	Data    []ConversationData `json:"data"`
}

type ConversationData struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Inputs       map[string]interface{} `json:"inputs"`
	Status       string                 `json:"status"`
	Introduction string                 `json:"introduction"`
	CreatedAt    int64                  `json:"created_at"`
}

func HandleGetConversation(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := ReadResponseBody(resp)
	if err != nil {
		return err
	}
	return HandleConversationData(resp, body, db, user)
}

func HandleConversationData(resp *http.Response, body []byte, db *gorm.DB, user models.User) error {

	if resp == nil || db == nil || user.ID.String() == "" {
		return fmt.Errorf("invalid input: nil or uninitialized parameter")
	}

	var originalResponse ConversationResponse
	if err := json.Unmarshal(body, &originalResponse); err != nil {
		logging.Warning("Failed to decode incoming data: " + err.Error())
		return err
	}

	filteredConversations := make([]ConversationData, 0)

	for _, conversation := range originalResponse.Data {
		conversationID := conversation.ID
		if conversationID == "" {
			logging.Warning("Invalid or missing conversation ID in incoming data")
			continue
		}

		dbConversation, err := models.GetConversationByID(conversationID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// filteredConversations = append(filteredConversations, conversation)
			} else {
				logging.Warning("Error checking conversation existence: " + err.Error())
				return err
			}
		} else {
			if dbConversation.Owner == user.ID.String() {
				filteredConversations = append(filteredConversations, conversation)
			}
		}
	}

	originalResponse.Data = filteredConversations

	modifiedBody, err := json.Marshal(originalResponse)
	if err != nil {
		return err
	}

	WriteResponseBody(resp, modifiedBody)
	return nil
}
