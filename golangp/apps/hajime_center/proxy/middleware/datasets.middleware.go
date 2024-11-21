package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"net/http"
)

func HandlePostDatasetInit(resp *http.Response, r *http.Request, db *gorm.DB, user models.User) error {
	body, err := ReadResponseBody(resp)
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

	if responseData.Dataset.ID == "" {
		return errors.New("create dataset error")
	}

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
	datasetData, err := StructToMap(createdDataset)
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

	WriteResponseBody(resp, modifiedBody)
	return nil
}

func HandleGetAllDatasets(resp *http.Response, r *http.Request, user models.User) error {
	body, err := ReadResponseBody(resp)
	if err != nil {
		return err
	}
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
		if err := MapToStruct(incomingAppData, &dataset); err != nil {
			logging.Warning("Failed to convert incoming app data to struct: " + err.Error())
			continue
		}

		// Fetch the app from the database
		dbApp, err := models.GetDatasetByID(id)
		dbAppData, err := StructToMap(dbApp)
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
	WriteResponseBody(resp, modifiedBody)
	return nil
}

func HandleDeleteDatasets(resp *http.Response, r *http.Request, user models.User) error {
	vars := mux.Vars(r)
	datasetID, ok := vars["dataset_id"]
	if !ok {
		logging.Warning("App ID is missing in the request URL")
		return fmt.Errorf("app ID is required")
	}

	if err := models.DeleteDatasetByID(datasetID); err != nil {
		logging.Warning("Failed to delete app: " + err.Error())
		return err
	}

	resp.StatusCode = http.StatusNoContent
	WriteResponseBody(resp, []byte(""))
	return nil
}
