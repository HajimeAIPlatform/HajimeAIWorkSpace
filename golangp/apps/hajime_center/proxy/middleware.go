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
	"strings"
)

// OriginalResponse defines the structure for the original response
type OriginalResponse struct {
	Data    []map[string]interface{} `json:"data"`
	Limit   int                      `json:"limit"`
	Total   int                      `json:"total"`
	Page    int                      `json:"page"`
	HasMore bool                     `json:"has_more"`
}

func ModifyResponse(w *http.Response, r *http.Request) error {
	if strings.HasPrefix(r.URL.Path, "/console/api/apps") {
		db := initializers.DB
		switch r.Method {
		case http.MethodGet:
			return handleGetRequest(w, r, db)
		case http.MethodPost:
			return handlePostRequest(w, r, db)
		case http.MethodPut:
			return handlePutRequest(w, r, db)
		case http.MethodDelete:
			return handleDeleteRequest(w, r, db)
		default:
			return nil
		}
	}
	return nil
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

func handleGetRequest(resp *http.Response, r *http.Request, db *gorm.DB) error {
	vars := mux.Vars(r)
	appID := vars["app_id"]

	body, err := readResponseBody(resp)
	if err != nil {
		return err
	}

	if appID != "" {
		return handleGetSingleApp(resp, body, appID, db)
	}

	return handleGetAllApps(resp, body, db)
}

func handleGetSingleApp(resp *http.Response, body []byte, appID string, db *gorm.DB) error {
	app, err := models.GetHajimeAppByID(db, appID)
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

func handleGetAllApps(resp *http.Response, body []byte, db *gorm.DB) error {
	var originalResponse OriginalResponse
	if err := json.Unmarshal(body, &originalResponse); err != nil {
		logging.Warning("Failed to decode incoming data: " + err.Error())
		return err
	}

	for i, incomingAppData := range originalResponse.Data {
		id, ok := incomingAppData["id"].(string)
		if !ok {
			logging.Warning("Invalid or missing ID in incoming app data")
			continue
		}

		var incomingApp models.HajimeApps
		if err := mapToStruct(incomingAppData, &incomingApp); err != nil {
			logging.Warning("Failed to convert incoming app data to struct: " + err.Error())
			return err
		}

		dbApp, err := models.GetHajimeAppByID(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := models.CreateHajimeApp(db, incomingApp); err != nil {
					logging.Warning("Failed to create app: " + err.Error())
					return err
				}
				fmt.Println("App added:", id)
			} else {
				logging.Warning("Error checking app existence: " + err.Error())
				return err
			}
		} else {
			dbAppData, err := structToMap(dbApp)
			if err != nil {
				logging.Warning("Failed to convert db app to map: " + err.Error())
				return err
			}

			for key, value := range dbAppData {
				incomingAppData[key] = value
			}

			originalResponse.Data[i] = incomingAppData
		}
	}

	modifiedBody, err := json.Marshal(originalResponse)
	if err != nil {
		logging.Warning("Failed to encode response: " + err.Error())
		return err
	}

	writeResponseBody(resp, modifiedBody)
	return nil
}

func handlePostRequest(resp *http.Response, r *http.Request, db *gorm.DB) error {
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

	if err := models.CreateHajimeApp(db, app); err != nil {
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

func handlePutRequest(resp *http.Response, r *http.Request, db *gorm.DB) error {
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
			if err := models.CreateHajimeApp(db, app); err != nil {
				logging.Warning("Failed to create app: " + err.Error())
				return err
			}
			fmt.Println("App created:", app.ID)
		} else {
			logging.Warning("Failed to find app: " + err.Error())
			return err
		}
	} else {
		if err := models.UpdateHajimeApp(db, app); err != nil {
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

func handleDeleteRequest(resp *http.Response, r *http.Request, db *gorm.DB) error {
	vars := mux.Vars(r)
	appID, ok := vars["app_id"]
	if !ok {
		logging.Warning("App ID is missing in the request URL")
		return fmt.Errorf("app ID is required")
	}

	if err := models.DeleteHajimeApp(db, appID); err != nil {
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
