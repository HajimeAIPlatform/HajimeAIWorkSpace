package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"net/http"
)

func HandleInstallGetRequest(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := ReadResponseBody(resp)
	if err != nil {
		return err
	}
	return HandleInstallGetAllApps(resp, body, db, user)
}

func HandleInstallGetAllApps(resp *http.Response, body []byte, db *gorm.DB, user models.User) error {
	var originalResponse InstallAppResponse
	if err := json.Unmarshal(body, &originalResponse); err != nil {
		logging.Warning("Failed to decode incoming data: " + err.Error())
		return err
	}

	for _, incomingAppData := range originalResponse.InstalledApps {
		appID := incomingAppData.App.ID
		if appID == "" {
			logging.Warning("Invalid or missing app ID in incoming app data")
			continue
		}

		InstallAppID := incomingAppData.ID
		if InstallAppID == "" {
			logging.Warning("Invalid or missing install app ID in incoming app data")
			continue
		}

		var incomingApp models.HajimeApps
		if err := MapToStructApps(incomingAppData.App, &incomingApp); err != nil {
			logging.Warning("Failed to convert incoming app data to struct: " + err.Error())
			return err
		}

		dbApp, err := models.GetHajimeAppByID(appID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new app entry
				incomingApp.InstallAppID = InstallAppID
				if err := models.CreateHajimeApp(incomingApp); err != nil {
					logging.Warning("Failed to create app: " + err.Error())
					return err
				}
				fmt.Println("App added:", appID)
			} else {
				logging.Warning("Error checking app existence: " + err.Error())
				return err
			}
		} else {
			// Update existing app entry with InstallAppID
			dbApp.InstallAppID = InstallAppID
			if err := db.Save(&dbApp).Error; err != nil {
				logging.Warning("Failed to update app: " + err.Error())
				return err
			}
		}
	}

	modifiedBody, err := json.Marshal(originalResponse)
	if err != nil {
		return err
	}

	WriteResponseBody(resp, modifiedBody)
	return nil
}
