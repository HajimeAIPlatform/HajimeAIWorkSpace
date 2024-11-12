package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"io/ioutil"
	"net/http"
	"regexp"
)

// OriginalResponse defines the structure for the original response
type OriginalResponse struct {
	Data    []map[string]interface{} `json:"data"`
	Limit   int                      `json:"limit"`
	Total   int                      `json:"total"`
	Page    int                      `json:"page"`
	HasMore bool                     `json:"has_more"`
}

type InstallAppResponse struct {
	InstalledApps []InstalledApps `json:"installed_apps"`
}

type InstalledApps struct {
	ID               string            `json:"id"`
	App              models.HajimeApps `json:"app"`
	AppOwnerTenantID string            `json:"app_owner_tenant_id"`
	IsPinned         bool              `json:"is_pinned"`
	LastUsedAt       int64             `json:"last_used_at"`
	Editable         bool              `json:"editable"`
	Uninstallable    bool              `json:"uninstallable"`
}

type RecommendedAPP struct {
	App struct {
		Icon           string `json:"icon"`
		IconBackground string `json:"icon_background"`
		ID             string `json:"id"`
		Mode           string `json:"mode"`
		Name           string `json:"name"`
	} `json:"app"`
	ID               string      `json:"id"`
	AppID            string      `json:"app_id"`
	Category         string      `json:"category"`
	Copyright        interface{} `json:"copyright"`
	CustomDisclaimer interface{} `json:"custom_disclaimer"`
	Description      interface{} `json:"description"`
	IsListed         bool        `json:"is_listed"`
	Position         int64       `json:"position"`
	PrivacyPolicy    interface{} `json:"privacy_policy"`
}

type NoAuthApp struct {
	Categories     []string         `json:"categories"`
	RecommendedAPP []RecommendedAPP `json:"recommended_apps"`
}

func ModifyResponse(w *http.Response, r *http.Request, user models.User) error {

	db := initializers.DB
	if r.URL.Path == "/console/api/apps" || isAppIDPath(r.URL.Path) {
		switch r.Method {
		case http.MethodGet:
			return handleGetRequest(w, r, db, user)
		case http.MethodPost:
			return handlePostRequest(w, r, db, user)
		case http.MethodPut:
			return handlePutRequest(w, r, db, user)
		case http.MethodDelete:
			return handleDeleteRequest(w, r, db, user)
		default:
			return nil
		}
	}

	if r.URL.Path == "/console/api/installed-apps" {
		switch r.Method {
		case http.MethodGet:
			return handleInstallGetRequest(w, r, db, user)
		default:
			return nil
		}
	}
	return nil
}

func isAppIDPath(path string) bool {
	// 匹配 "/console/api/apps/{app_id}"，确保后面没有其他路径
	matched, _ := regexp.MatchString(`^/console/api/apps/[a-fA-F0-9\-]+/?$`, path)
	return matched
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	resp.Body.Close()
	return buf.Bytes(), err
}

func writeResponseBody(resp *http.Response, body []byte) {
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))
	resp.Header.Set("Content-Length", fmt.Sprint(len(body)))
}

func handleInstallGetRequest(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := readResponseBody(resp)
	if err != nil {
		return err
	}
	return handleInstallGetAllApps(resp, body, db, user)
}

func handleInstallGetAllApps(resp *http.Response, body []byte, db *gorm.DB, user models.User) error {
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
		if err := mapToStructApps(incomingAppData.App, &incomingApp); err != nil {
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

	writeResponseBody(resp, modifiedBody)
	return nil
}

func handleGetRequest(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	vars := mux.Vars(r)
	appID := vars["app_id"]

	body, err := readResponseBody(resp)
	if err != nil {
		return err
	}

	if appID != "" {
		return handleGetSingleApp(resp, body, appID)
	}

	return handleGetAllApps(resp, r, body, user)
}

func handleGetSingleApp(resp *http.Response, body []byte, appID string) error {
	app, err := models.GetHajimeAppByID(appID)
	if err != nil {
		logging.Warning("Failed to fetch app: " + err.Error())
		return err
	}

	var originalData map[string]interface{}
	if err := json.Unmarshal(body, &originalData); err != nil {
		return err
	}

	appData, err := structToMap(app)
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

	writeResponseBody(resp, modifiedBody)
	return nil
}

func handleGetAllApps(resp *http.Response, r *http.Request, body []byte, user models.User) error {
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
		if err := mapToStruct(incomingAppData, &incomingApp); err != nil {
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
			dbAppData, err := structToMap(dbApp)
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

	writeResponseBody(resp, modifiedBody)
	return nil
}

func handlePostRequest(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := readResponseBody(resp)
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
		writeResponseBody(resp, body)
		return nil
	}

	var app models.HajimeApps
	if err := json.Unmarshal(body, &app); err != nil {
		logging.Warning("Failed to parse response body: " + err.Error())
		return err
	}
	app.Owner = user.ID.String()

	fmt.Println("user", app.Owner, user)

	if err := models.CreateHajimeApp(app); err != nil {
		logging.Warning("Failed to create app: " + err.Error())
		return err
	}

	var createdApp models.HajimeApps
	if err := db.First(&createdApp, "id = ?", app.ID).Error; err != nil {
		logging.Warning("Failed to retrieve created app: " + err.Error())
		return err
	}

	appData, err := structToMap(createdApp)
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

	writeResponseBody(resp, modifiedBody)
	return nil
}

func handlePutRequest(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := readResponseBody(resp)
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

	appData, err := structToMap(updatedApp)
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

	writeResponseBody(resp, modifiedBody)
	return nil
}

func handleDeleteRequest(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
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

	resp.StatusCode = http.StatusNoContent
	writeResponseBody(resp, []byte(""))
	return nil
}

func structToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func mapToStruct(data map[string]interface{}, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, result)
}

func mapToStructApps(data models.HajimeApps, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, result)
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

	// Update the app's isPublish status
	existingApp.IsPublish = true // or whatever logic you need
	if err := models.UpdateHajimeApp(existingApp); err != nil {
		logging.Warning("Failed to update app: " + err.Error())
		http.Error(w, "Failed to update app", http.StatusInternalServerError)
		return
	}
	// 设置响应头为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 返回成功状态和 JSON 响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
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

	// Iterate over apps and build the RecommendedAPP list
	for index, app := range apps {
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

func HandlePostDatasetInit(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := readResponseBody(resp)
	if err != nil {
		logging.Warning("Failed to read response body: " + err.Error())
		return err
	}

	var responseData struct {
		Batch     string         `json:"batch"`
		Dataset   models.Dataset `json:"dataset"`
		Documents []interface{}  `json:"documents"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		logging.Warning("Failed to parse response body: " + err.Error())
		return err
	}

	// 设置数据集的所有者
	responseData.Dataset.Owner = user.ID.String()

	// 保存数据集到数据库
	if err := models.SaveDataset(&responseData.Dataset); err != nil {
		logging.Warning("Failed to save dataset: " + err.Error())
		return err
	}

	// 从数据库中检索已创建的数据集
	var createdDataset models.Dataset
	if err := db.First(&createdDataset, "id = ?", responseData.Dataset.ID).Error; err != nil {
		logging.Warning("Failed to retrieve created dataset: " + err.Error())
		return err
	}

	// 将数据集转换为 map
	datasetData, err := structToMap(createdDataset)
	if err != nil {
		return err
	}

	// 更新原始数据
	originalData := make(map[string]interface{})
	if err := json.Unmarshal(body, &originalData); err != nil {
		return err
	}

	originalData["dataset"] = datasetData

	// 序列化修改后的数据
	modifiedBody, err := json.Marshal(originalData)
	if err != nil {
		return err
	}

	writeResponseBody(resp, modifiedBody)
	return nil
}

func handleGetAllDatasets(resp *http.Response, r *http.Request, body []byte, user models.User) error {
	var originalResponse OriginalResponse
	if err := json.Unmarshal(body, &originalResponse); err != nil {
		logging.Warning("Failed to decode incoming data: " + err.Error())
		return err
	}

	filteredData := []map[string]interface{}{}

	for _, incomingAppData := range originalResponse.Data {
		id, ok := incomingAppData["id"].(string)
		if !ok {
			logging.Warning("Invalid or missing ID in incoming app data")
			continue
		}

		var dataset models.Dataset
		if err := mapToStruct(incomingAppData, &dataset); err != nil {
			logging.Warning("Failed to convert incoming app data to struct: " + err.Error())
			continue
		}

		// Fetch the app from the database
		dbApp, err := models.GetDatasetByID(id)
		dbAppData, err := structToMap(dbApp)
		if err != nil {
			logging.Warning("Failed to convert db app to map: " + err.Error())
			continue
		}

		// Merge database data into incoming app data
		for key, value := range dbAppData {
			incomingAppData[key] = value
		}

		// Filter based on user role and ownership
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

	// Update the original response with the filtered data
	originalResponse.Data = filteredData

	// Encode the modified response
	modifiedBody, err := json.Marshal(originalResponse)
	if err != nil {
		logging.Warning("Failed to encode response: " + err.Error())
		return err
	}

	// Write the response
	writeResponseBody(resp, modifiedBody)
	return nil
}
