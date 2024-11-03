package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Apps struct {
	ID           string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	PrePrompt    string    `gorm:"type:text" json:"pre_prompt"`
	Type         string    `gorm:"type:varchar(255);not null" json:"type"`
	RoleIndustry string    `gorm:"type:varchar(255)" json:"role_industry,omitempty"`
	RoleSettings string    `gorm:"type:varchar(255)" json:"role_settings,omitempty"`
	Model        string    `gorm:"type:varchar(255);not null" json:"model"`
	DatasetId    string    `gorm:"type:text" json:"dataset_id,omitempty"`
	IsPublished  bool      `gorm:"type:boolean;default:false" json:"is_published"`
	CreateBy     uuid.UUID `gorm:"type:uuid;default:null" json:"create_by"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
}

// BeforeSave 钩子函数，在保存之前调用，将 ModelConfig 转换为 JSON 字符串
func (app *Apps) BeforeSave(tx *gorm.DB) error {
	if app.Type != "assistant" { //assistant | knowledge
		app.RoleIndustry = ""
		app.RoleSettings = ""
	}
	return nil
}

type CreateAppsInput struct {
	Name         string   `json:"name,required"`
	Type         string   `json:"type,required"`
	Description  string   `json:"description"`
	Model        string   `json:"model,required"`
	PrePrompt    string   `json:"pre_prompt"`
	RoleIndustry string   `json:"role_industry,omitempty"`
	RoleSettings string   `json:"role_settings,omitempty"`
	DatasetId    []string `json:"dataset_id,omitempty"` // 修改为 []string
}

type UpdateAppsInput struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,required"`
	Description  string   `json:"description,omitempty"`
	PrePrompt    string   `json:"pre_prompt,omitempty"`
	RoleIndustry string   `json:"role_industry,omitempty"`
	RoleSettings string   `json:"role_settings,omitempty"`
	Model        string   `json:"model,omitempty"`
	DatasetId    []string `json:"dataset_id,omitempty"` // 修改为 []string
}

type GetAppsListInputViaType struct {
	Type string `json:"type,omitempty"`
}
