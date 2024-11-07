package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hajime/golangp/common/logging"
)

type HajimeApps struct {
	ID               string `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID         int    `gorm:"not null" json:"tenant_id"`
	Mode             string `gorm:"type:varchar(50)" json:"mode"`
	Name             string `gorm:"type:varchar(100)" json:"name"`
	Description      string `gorm:"type:text" json:"description"`
	AppModelConfigID int    `gorm:"not null" json:"app_model_config_id"`
	WorkflowID       int    `gorm:"not null" json:"workflow_id"`
	Status           string `gorm:"type:varchar(50)" json:"status"`
	Owner            string `gorm:"type:varchar(100) default:''" json:"owner"`
	IsPublic         bool   `gorm:"not null default:false" json:"is_public"`
}

// CreateHajimeApp 创建一个新的 HajimeApps
func CreateHajimeApp(db *gorm.DB, app HajimeApps) error {
	if err := db.Create(&app).Error; err != nil {
		return err
	}
	fmt.Println("App created:", app)
	return nil
}

// GetHajimeAppByID 根据ID获取HajimeApps
func GetHajimeAppByID(db *gorm.DB, id string) (HajimeApps, error) {
	var app HajimeApps
	if err := db.First(&app, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return HajimeApps{}, err
		}
		return HajimeApps{}, err
	}
	return app, nil
}

// UpdateHajimeApp 更新HajimeApps
func UpdateHajimeApp(db *gorm.DB, app HajimeApps) error {
	// Find the existing app by ID
	var existingApp HajimeApps
	if err := db.First(&existingApp, "id = ?", app.ID).Error; err != nil {
		logging.Warning("Failed to find app: " + err.Error())
		return err
	}

	// Update fields
	existingApp.Name = app.Name
	existingApp.Description = app.Description

	// Save changes
	if err := db.Save(&existingApp).Error; err != nil {
		logging.Warning("Failed to update app: " + err.Error())
		return err
	}

	return nil
}

// DeleteHajimeApp 删除HajimeApps
func DeleteHajimeApp(db *gorm.DB, id string) error {
	if err := db.Delete(&HajimeApps{}, "id = ?", id).Error; err != nil {
		return err
	}
	fmt.Println("App deleted with ID:", id)
	return nil
}

// GetAllHajimeApps 获取所有 HajimeApps
func GetAllHajimeApps(db *gorm.DB) ([]HajimeApps, error) {
	var apps []HajimeApps
	if err := db.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}
