package main

import (
	openai_config "HajimeAIWorkSpace/common/apps/hajime_center/chat-config"
	"HajimeAIWorkSpace/common/apps/hajime_center/constants"
	"HajimeAIWorkSpace/common/apps/hajime_center/initializers"
	"HajimeAIWorkSpace/common/apps/hajime_center/logger"
	"HajimeAIWorkSpace/common/apps/hajime_center/models"
	"HajimeAIWorkSpace/common/apps/hajime_center/utils"
	"fmt"
	"github.com/alecthomas/kong"
	"gorm.io/gorm"
	"strings"
	"time"
)

func init() {
	config, err := initializers.LoadEnv(".")
	if err != nil {
		logger.Danger("🚀 Could not load environment variables %s", err.Error())
	}

	initializers.ConnectDB(&config)
}

func removeAllAdmins(DB *gorm.DB) {
	var adminUsers []models.User
	res := DB.Find(&adminUsers, "role = ?", constants.RoleAdmin)
	if res.Error != nil {
		logger.Warning("Error finding admin users %s", res.Error.Error())
	}

	for _, adminUser := range adminUsers {
		userToDelete := adminUser.Email
		res := DB.Delete(&adminUser)
		if res.Error != nil {
			logger.Warning("Error deleting admin user %s", res.Error.Error())
		}
		logger.Info("Previous admin user %s deleted successfully", userToDelete)
	}
}

func SetupAdmin(DB *gorm.DB) {
	kong.Parse(&openai_config.CLI)
	OpenaiConfig := openai_config.LoadChatConfig()
	adminPassword := OpenaiConfig.AdminPassword

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		logger.Danger("Error hashing password %s", err.Error())
	}

	removeAllAdmins(DB)

	for _, adminEmail := range OpenaiConfig.AdminEmail {
		now := time.Now()
		newUser := models.User{
			Name:      "Admin Admin",
			Email:     strings.ToLower(adminEmail),
			Password:  hashedPassword,
			Role:      constants.RoleAdmin,
			Verified:  true,
			Photo:     "test",
			Provider:  "local",
			CreatedAt: now,
			UpdatedAt: now,
		}

		var adminUser models.User
		res := DB.First(&adminUser, "email = ?", adminEmail)
		if res.Error != nil {
			logger.Info("Admin user %s does not exist, creating one", adminEmail)
		} else {
			res := DB.Delete(&adminUser)
			if res.Error != nil {
				logger.Warning("Error deleting exist admin user %s", res.Error.Error())
			}
			logger.Info("Existing Admin user deleted successfully")
		}

		result := DB.Create(&newUser)

		if result.Error != nil && strings.Contains(result.Error.Error(), "duplicated key not allowed") {
			logger.Warning("Admin email already exists")
			return
		} else if result.Error != nil {
			logger.Danger("Error creating admin user", result.Error)
		}

		logger.Info("Admin user %s created successfully", adminEmail)
	}
}

func main() {
	// Ensure the uuid-ossp extension is created
	if err := initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		logger.Danger("🚀 Could not create uuid-ossp extension: %v", err)
	}

	err := initializers.DB.AutoMigrate(&models.User{}, &models.BillingHistory{}, &models.Post{},&models.Apps{},&models.Dataset{},&models.Document{},&models.UploadFile{},&models.Conversation{},&models.Message{},&models.MessageFile{})

	if err != nil {
		logger.Danger(fmt.Sprintf("🚀 Could not migrate User model: %v", err))
	}

	SetupAdmin(initializers.DB)
	fmt.Println("👍 Migration complete")
}
