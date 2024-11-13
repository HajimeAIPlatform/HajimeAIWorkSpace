package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"net/http"
)

func HandleGetApp(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	vars := mux.Vars(r)
	appID := vars["app_id"]

	body, err := ReadResponseBody(resp)
	if err != nil {
		return err
	}

	if appID != "" {
		return HandleGetSingleApp(resp, body, appID)
	}

	return HandleGetAllApps(resp, r, body, user)
}

func HandleGetSingleApp(resp *http.Response, body []byte, appID string) error {
	app, err := models.GetHajimeAppByID(appID)
	if err != nil {
		logging.Warning("Failed to fetch app: " + err.Error())
		return err
	}

	var originalData map[string]interface{}
	if err := json.Unmarshal(body, &originalData); err != nil {
		return err
	}

	appData, err := StructToMap(app)
	if err != nil {
		return err
	}

	for key, value := range appData {
		originalData[key] = value
	}

	modifiedBody, err := json.Marshal(originalData)
	if err != nil {
		return err
	}

	WriteResponseBody(resp, modifiedBody)
	return nil
}

func HandleGetAllApps(resp *http.Response, r *http.Request, body []byte, user models.User) error {
	var originalResponse OriginalResponse
	if err := json.Unmarshal(body, &originalResponse); err != nil {
		logging.Warning("Failed to decode incoming data: " + err.Error())
		return err
	}

	installedApps, err := FetchInstalledApps(r)
	if err != nil {
		logging.Warning("Failed to fetch installed apps: " + err.Error())
		return err
	}

	filteredData := []map[string]interface{}{}

	for _, incomingAppData := range originalResponse.Data {
		id, ok := incomingAppData["id"].(string)
		if !ok {
			logging.Warning("Invalid or missing ID in incoming app data")
			continue
		}

		var incomingApp models.HajimeApps
		if err := MapToStruct(incomingAppData, &incomingApp); err != nil {
			logging.Warning("Failed to convert incoming app data to struct: " + err.Error())
			continue
		}

		for _, installedApp := range installedApps {
			if installedApp.App.ID == id {
				incomingApp.InstallAppID = installedApp.ID
				break
			}
		}

		dbApp, err := models.GetHajimeAppByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := models.CreateHajimeApp(incomingApp); err != nil {
					logging.Warning("Failed to create app: " + err.Error())
					continue
				}
			} else {
				logging.Warning("Error checking app existence: " + err.Error())
				continue
			}
		} else {
			dbAppData, err := StructToMap(dbApp)
			if err != nil {
				logging.Warning("Failed to convert db app to map: " + err.Error())
				continue
			}

			for key, value := range dbAppData {
				incomingAppData[key] = value
			}
		}

		owner, ok := incomingAppData["owner"].(string)
		if user.Role == "admin" {
			if owner == "" || owner == user.ID.String() {
				filteredData = append(filteredData, incomingAppData)
			}
		} else {
			if ok && owner == user.ID.String() {
				filteredData = append(filteredData, incomingAppData)
			}
		}
	}

	originalResponse.Data = filteredData

	modifiedBody, err := json.Marshal(originalResponse)
	if err != nil {
		logging.Warning("Failed to encode response: " + err.Error())
		return err
	}

	WriteResponseBody(resp, modifiedBody)
	return nil
}

func HandlePostApp(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := ReadResponseBody(resp)
	if err != nil {
		logging.Warning("Failed to read response body: " + err.Error())
		return err
	}

	var originalData map[string]interface{}
	if err := json.Unmarshal(body, &originalData); err != nil {
		logging.Warning("Failed to parse original response body: " + err.Error())
		return err
	}

	if code, exists := originalData["code"]; exists && code == "bad_request" {
		WriteResponseBody(resp, body)
		return nil
	}

	var app models.HajimeApps
	if err := json.Unmarshal(body, &app); err != nil {
		logging.Warning("Failed to parse response body: " + err.Error())
		return err
	}
	app.Owner = user.ID.String()

	if err := models.CreateHajimeApp(app); err != nil {
		logging.Warning("Failed to create app: " + err.Error())
		return err
	}

	err = models.UpdateAppAmount(user.ID.String(), 1)
	if err != nil {
		logging.Warning("Failed to retrieve created app:: " + err.Error())
		return err
	}

	var createdApp models.HajimeApps
	if err := db.First(&createdApp, "id = ?", app.ID).Error; err != nil {
		logging.Warning("Failed to retrieve created app: " + err.Error())
		return err
	}

	appData, err := StructToMap(createdApp)
	if err != nil {
		return err
	}

	for key, value := range appData {
		originalData[key] = value
	}

	modifiedBody, err := json.Marshal(originalData)
	if err != nil {
		return err
	}

	WriteResponseBody(resp, modifiedBody)
	return nil
}

func HandlePutApp(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := ReadResponseBody(resp)
	if err != nil {
		logging.Warning("Failed to read response body: " + err.Error())
		return err
	}

	var originalData map[string]interface{}
	if err := json.Unmarshal(body, &originalData); err != nil {
		logging.Warning("Failed to parse original response body: " + err.Error())
		return err
	}

	var app models.HajimeApps
	if err := json.Unmarshal(body, &app); err != nil {
		logging.Warning("Failed to parse response body: " + err.Error())
		return err
	}

	var existingApp models.HajimeApps
	if err := db.Where("id = ?", app.ID).First(&existingApp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.Owner = user.ID.String()
			if err := models.CreateHajimeApp(app); err != nil {
				logging.Warning("Failed to create app: " + err.Error())
				return err
			}
			fmt.Println("App created:", app.ID)
		} else {
			logging.Warning("Failed to find app: " + err.Error())
			return err
		}
	} else {
		if err := models.UpdateHajimeApp(app); err != nil {
			logging.Warning("Failed to update app: " + err.Error())
			return err
		}
	}

	var updatedApp models.HajimeApps
	if err := db.Where("id = ?", app.ID).First(&updatedApp).Error; err != nil {
		logging.Warning("Failed to retrieve updated app: " + err.Error())
		return err
	}

	appData, err := StructToMap(updatedApp)
	if err != nil {
		logging.Warning("Failed to convert app to map: " + err.Error())
		return err
	}

	for key, value := range appData {
		originalData[key] = value
	}

	modifiedBody, err := json.Marshal(originalData)
	if err != nil {
		logging.Warning("Failed to marshal modified response: " + err.Error())
		return err
	}

	WriteResponseBody(resp, modifiedBody)
	return nil
}

func HandleDeleteApp(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	vars := mux.Vars(r)
	appID, ok := vars["app_id"]
	if !ok {
		logging.Warning("App ID is missing in the request URL")
		return fmt.Errorf("app ID is required")
	}

	if err := models.DeleteHajimeApp(appID); err != nil {
		logging.Warning("Failed to delete app: " + err.Error())
		return err
	}
	err := models.UpdateAppAmount(user.ID.String(), -1)
	if err != nil {
		logging.Warning("Failed to retrieve created app:: " + err.Error())
		return err
	}

	resp.StatusCode = http.StatusNoContent
	WriteResponseBody(resp, []byte(""))
	return nil
}

func GetAllNoAuthApp(w http.ResponseWriter, r *http.Request) {

	// Call GetAllHajimeAppsNoAuth to get published apps
	apps, err := models.GetAllHajimeAppsNoAuth()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize a single NoAuthApp with empty categories
	noAuthApp := NoAuthApp{
		Categories:     []string{},
		RecommendedAPP: []RecommendedAPP{},
	}

	installedApps, err := FetchInstalledApps(r)
	if err != nil {
		logging.Warning("Failed to fetch installed apps: " + err.Error())
		return
	}
	installedAppMap := make(map[string]string)
	for _, installedApp := range installedApps {
		installedAppMap[installedApp.App.ID] = installedApp.ID
	}

	// Iterate over apps and build the RecommendedAPP list
	for index, app := range apps {
		if id, exists := installedAppMap[app.ID]; exists {
			app.InstallAppID = id
		}
		recommendedApp := RecommendedAPP{
			App: struct {
				Icon           string `json:"icon"`
				IconBackground string `json:"icon_background"`
				ID             string `json:"id"`
				Mode           string `json:"mode"`
				Name           string `json:"name"`
			}{
				Icon:           app.Icon,
				IconBackground: app.IconBackground,
				ID:             app.ID,
				Mode:           app.Mode,
				Name:           app.Name,
			},
			ID:               app.InstallAppID,
			AppID:            app.ID,
			Category:         "", // Default to empty string
			Copyright:        nil,
			CustomDisclaimer: nil,
			Description:      nil,
			IsListed:         false,
			Position:         int64(index), // Use index as Position
			PrivacyPolicy:    nil,
		}

		// Append each app to the RecommendedAPP slice
		noAuthApp.RecommendedAPP = append(noAuthApp.RecommendedAPP, recommendedApp)
	}

	// If no apps found, ensure RecommendedAPP is an empty slice
	if len(apps) == 0 {
		noAuthApp.RecommendedAPP = []RecommendedAPP{}
	}

	// Convert the single NoAuthApp to JSON
	response, err := json.Marshal(noAuthApp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header to JSON format
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// HandlePublish is a custom handler for the /publish route
func HandlePublish(w http.ResponseWriter, r *http.Request) {
	// Extract the app_id from the URL
	vars := mux.Vars(r) // Assuming you're using Gorilla Mux
	appID := vars["app_id"]

	// Define the app structure
	var existingApp models.HajimeApps // Replace 'App' with your actual model struct
	db := initializers.DB
	// Check if the app exists in the database
	if err := db.Where("id = ?", appID).First(&existingApp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "App not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	isGreater, err := models.IsAppPublishAmountGreaterThanTen(existingApp.Owner)
	if err != nil {
		logging.Warning("This app is not your own: " + err.Error())
		http.Error(w, "This app is not your own", http.StatusInternalServerError)
		return
	}

	if isGreater {
		logging.Warning("You have exceeded the maximum number of publishes, which is 10.")
		http.Error(w, "You have exceeded the maximum number of publishes, which is 10.", http.StatusForbidden)
		return
	}

	// Update the app's isPublish status
	existingApp.IsPublish = true // or whatever logic you need
	if err := models.UpdateHajimeApp(existingApp); err != nil {
		logging.Warning("Failed to update app: " + err.Error())
		http.Error(w, "Failed to update app", http.StatusInternalServerError)
		return
	}
	err = models.UpdateAppPublishAmount(existingApp.Owner, 1)
	if err != nil {
		logging.Warning("Failed to update app: " + err.Error())
		http.Error(w, "Failed to update app", http.StatusInternalServerError)
		return
	}

	// 设置响应头为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 返回成功状态和 JSON 响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "Publish success"})
}

func HandleUnpublish(w http.ResponseWriter, r *http.Request) {
	// Extract the app_id from the URL
	vars := mux.Vars(r) // Assuming you're using Gorilla Mux
	appID := vars["app_id"]

	// Define the app structure
	var existingApp models.HajimeApps // Replace 'HajimeApps' with your actual model struct
	db := initializers.DB

	// Check if the app exists in the database
	if err := db.Where("id = ?", appID).First(&existingApp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "App not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Update the app's isPublish status to false
	existingApp.IsPublish = false // or whatever logic you need
	if err := models.UpdateHajimeApp(existingApp); err != nil {
		logging.Warning("Failed to update app: " + err.Error())
		http.Error(w, "Failed to update app", http.StatusInternalServerError)
		return
	}
	err := models.UpdateAppPublishAmount(existingApp.Owner, -1)
	if err != nil {
		logging.Warning("Failed to update app: " + err.Error())
		http.Error(w, "Failed to update app", http.StatusInternalServerError)
		return
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Return success status and JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "Cancel publish success"})
}
