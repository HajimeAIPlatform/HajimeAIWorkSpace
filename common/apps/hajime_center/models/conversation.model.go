package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type Conversation struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AppID         uuid.UUID      `gorm:"type:uuid;not null;index:conversation_app_from_user_idx" json:"appID"`
	ModelProvider string         `gorm:"type:varchar(255)" json:"modelProvider,omitempty"`
	ModelName     string         `gorm:"type:varchar(255)" json:"modelName"`
	PrePrompt     string         `gorm:"type:text" json:"prePrompt,omitempty"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Inputs        datatypes.JSON `gorm:"type:jsonb" json:"inputs,omitempty"`
	FromAccountID uuid.UUID      `gorm:"type:uuid;index:message_account_idx" json:"fromAccountId"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
}

type Message struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AppID            uuid.UUID      `gorm:"type:uuid;not null;index:message_app_id_idx" json:"appID"`
	ModelProvider    string         `gorm:"type:varchar(255);default:'openai';not null" json:"modelProvider"`
	ModelName        string         `gorm:"type:varchar(255)" json:"modelName"`
	ConversationID   uuid.UUID      `gorm:"type:uuid;not null;index:message_conversation_id_idx" json:"conversationID"`
	Inputs           datatypes.JSON `gorm:"type:jsonb" json:"inputs,omitempty"`
	Query            string         `gorm:"type:text;not null" json:"query"`
	Message          datatypes.JSON `gorm:"type:jsonb;not null" json:"message,omitempty" `
	MessageTokens    int            `gorm:"default:0;not null" json:"messageTokens,omitempty"`
	MessageUnitPrice float64        `gorm:"type:numeric(10,4);not null" json:"messageUnitPrice,omitempty"`
	MessagePriceUnit float64        `gorm:"type:numeric(10,7);default:0.001;not null" json:"messagePriceUnit,omitempty"`
	Answer           string         `gorm:"type:text;not null" json:"answer"`
	AnswerTokens     int            `gorm:"default:0;not null" json:"answerTokens,omitempty"`
	AnswerUnitPrice  float64        `gorm:"type:numeric(10,4);not null" json:"answerUnitPrice,omitempty"`
	AnswerPriceUnit  float64        `gorm:"type:numeric(10,7);default:0.001;not null" json:"answerPriceUnit,omitempty"`
	TotalPrice       float64        `gorm:"type:numeric(10,7)" json:"totalPrice,omitempty"`
	Currency         string         `gorm:"type:varchar(255);not null" json:"currency,omitempty"`
	FromAccountID    uuid.UUID      `gorm:"type:uuid;index:message_account_idx" json:"fromAccountID,omitempty"`
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	Refs             datatypes.JSON `gorm:"type:jsonb" json:"refs,omitempty"`
	Cost             datatypes.JSON `gorm:"type:jsonb" json:"cost,omitempty"`
}

type MessageFile struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MessageID      uuid.UUID `gorm:"type:uuid;not null"`
	Type           string    `gorm:"type:varchar(255)"`
	URL            string    `gorm:"type:text"`
	TransferMethod string    `gorm:"type:varchar(255)"`
	UploadFileID   uuid.UUID `gorm:"type:uuid"`
	CreatedByRole  string    `gorm:"type:varchar(255);not null;default:'account'" json:"created_by_role"`
	CreatedBy      string    `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}
