package models

import (
	"hajime/golangp/apps/hajime_center/initializers"
)

type Conversation struct {
	ID           string   `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string   `gorm:"type:varchar(100)" json:"name"`
	Introduction string   `gorm:"type:text" json:"introduction"`
	Owner        string   `gorm:"type:varchar(100);default:''" json:"owner"`
	Status       string   `gorm:"type:varchar(50)" json:"status"`
	CreatedAt    UnixTime `gorm:"type:bigint" json:"created_at"`
}

func CreateConversation(id string, owner string) (*Conversation, error) {
	db := initializers.DB

	conversation := &Conversation{
		ID:           id,
		Name:         "New conversation",
		Owner:        owner,
		Introduction: "",
		Status:       "",
	}

	if err := db.Create(conversation).Error; err != nil {
		return nil, err
	}

	return conversation, nil
}

func UpdateConversation(conversation *Conversation) error {
	db := initializers.DB

	if err := db.Save(conversation).Error; err != nil {
		return err
	}

	return nil
}

func GetConversationByID(id string) (*Conversation, error) {
	db := initializers.DB

	var conversation Conversation

	if err := db.First(&conversation, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}
