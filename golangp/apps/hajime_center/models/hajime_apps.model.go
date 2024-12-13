package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/common/logging"
)

type UnixTime int64

type HajimeApps struct {
	ID               string   `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID         int      `gorm:"not null" json:"tenant_id"`
	Mode             string   `gorm:"type:varchar(50)" json:"mode"`
	Name             string   `gorm:"type:varchar(100)" json:"name"`
	Description      string   `gorm:"type:text" json:"description"`
	AppModelConfigID int      `gorm:"not null" json:"app_model_config_id"`
	WorkflowID       int      `gorm:"not null" json:"workflow_id"`
	Status           string   `gorm:"type:varchar(50)" json:"status"`
	Owner            string   `gorm:"type:varchar(100);default:''" json:"owner"`
	IsPublish        bool     `gorm:"not null;default:false" json:"is_publish"`
	Icon             string   `gorm:"type:varchar(100)" json:"icon"`
	IconBackground   string   `gorm:"type:varchar(100)" json:"icon_background"`
	CreatedAt        UnixTime `gorm:"type:bigint" json:"created_at"`
	PublishAt        UnixTime `gorm:"type:bigint" json:"publish_at"`
	InstallAppID     string   `gorm:"type:varchar(50)" json:"install_app_id"`
}

// CreateHajimeApp 创建一个新的 HajimeApps
func CreateHajimeApp(app HajimeApps) error {
	db := initializers.DB

	if err := db.Create(&app).Error; err != nil {
		return err
	}
	var user User
	ownerUUID, err := uuid.Parse(app.Owner)
	if err != nil {
		return fmt.Errorf("invalid UUID format for owner: %v", err)
	}
	if err := db.First(&user, ownerUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}
	if app.Mode == "workflow" {
		err := user.UpdateBalance(constants.CreateWorkflowPoints, "CreateWorkflowPoints")
		if err != nil {
			return err
		}
	} else {
		err := user.UpdateBalance(constants.CreateChatbotPoints, "CreateChatbotPoints")
		if err != nil {
			return err
		}
	}

	fmt.Println("App created:", app)
	return nil
}

// GetHajimeAppByID 根据ID获取HajimeApps
func GetHajimeAppByID(id string) (HajimeApps, error) {
	db := initializers.DB
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
func UpdateHajimeApp(app HajimeApps) error {
	// Find the existing app by ID
	db := initializers.DB
	var existingApp HajimeApps
	if err := db.First(&existingApp, "id = ?", app.ID).Error; err != nil {
		logging.Warning("Failed to find app: " + err.Error())
		return err
	}

	// Update fields
	existingApp.Name = app.Name
	existingApp.Description = app.Description
	existingApp.IsPublish = app.IsPublish

	// Save changes
	if err := db.Save(&existingApp).Error; err != nil {
		logging.Warning("Failed to update app: " + err.Error())
		return err
	}

	return nil
}

// DeleteHajimeApp 删除HajimeApps
func DeleteHajimeApp(id string) error {
	db := initializers.DB
	if err := db.Delete(&HajimeApps{}, "id = ?", id).Error; err != nil {
		return err
	}
	fmt.Println("App deleted with ID:", id)
	return nil
}

// GetAllHajimeApps 获取所有 HajimeApps
func GetAllHajimeApps() ([]HajimeApps, error) {
	db := initializers.DB
	var apps []HajimeApps
	if err := db.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func GetAllHajimeAppsNoAuth() ([]HajimeApps, error) {
	db := initializers.DB
	var apps []HajimeApps
	if err := db.Where("is_publish = ?", true).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

// GetUserByAppID 根据 AppID 获取应用的所有者
func GetUserByAppID(appID string) (User, error) {
	db := initializers.DB
	var app HajimeApps
	var user User

	// 查找应用
	if err := db.First(&app, "id = ?", appID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.First(&app, "install_app_id = ?", appID).Error; err != nil {
				return User{}, fmt.Errorf("app not found")
			}
		}
	}

	//TODO: Remove This
	app.Owner = ""

	// 检查 Owner 字段是否为空
	if app.Owner == "" {
		// 如果 Owner 为空，查找默认用户
		if err := db.First(&user, "email = ?", "hajime@gmail.com").Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return User{}, fmt.Errorf("default user not found")
			}
			return User{}, err
		}
		return user, nil
	}

	// 解析 Owner 字段为 UUID
	ownerUUID, err := uuid.Parse(app.Owner)
	if err != nil {
		return User{}, fmt.Errorf("invalid UUID format for owner: %v", err)
	}

	// 查找用户
	if err := db.First(&user, "id = ?", ownerUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, err
	}

	return user, nil
}
