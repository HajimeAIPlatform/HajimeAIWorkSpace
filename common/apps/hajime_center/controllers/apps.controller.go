package controllers

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/dify"
	"HajimeAIWorkSpace/common/apps/hajime_center/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AppsController struct {
	DB *gorm.DB
}

type AppsForIdResponse struct {
	dify.SiteStruct
	models.Apps
	DatasetId []string `json:"dataset_id,omitempty"`
	Datasets  []struct {
		Dataset dify.DatasetArray `json:"dataset"`
	} `json:"datasets"`
}

type AppsListResponse struct {
	models.Apps
	DatasetId []string `json:"dataset_id,omitempty"`
	Datasets  []struct {
		Dataset dify.DatasetArray `json:"dataset"`
	} `json:"datasets"`
}

func NewAppsController(DB *gorm.DB) AppsController {
	return AppsController{DB}
}

func (ac *AppsController) CreateApps(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateAppsInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("payload: ", payload.Description)

	difyClient := InitDifyClient()

	result, err := difyClient.CreateAppsHandler(payload.Name, payload.Description, payload.Model, payload.PrePrompt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to Create Apps: %v", err)})
		return
	}

	// Convert DatasetId to JSON string
	var datasetIdJSON string
	if len(payload.DatasetId) > 0 {
		datasetIdBytes, err := json.Marshal(payload.DatasetId)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "marshal dataset_id failed", "error": err})
			return
		}
		datasetIdJSON = string(datasetIdBytes)
	}

	apps := models.Apps{
		ID:           result.ID,
		Name:         result.Name,
		Description:  result.Description,
		Type:         payload.Type,
		RoleIndustry: payload.RoleIndustry,
		RoleSettings: payload.RoleSettings,
		PrePrompt:    payload.PrePrompt,
		Model:        payload.Model,
		DatasetId:    datasetIdJSON,
		CreateBy:     currentUser.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	resultDb := ac.DB.Create(&apps)
	if resultDb.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "create apps failed", "error": resultDb.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": apps})
}

func (ac *AppsController) UpdateApps(ctx *gin.Context) {
	var payload *models.UpdateAppsInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	difyClient := InitDifyClient()

	// Convert DatasetId to JSON string
	var datasetIdJSON string
	if len(payload.DatasetId) > 0 {
		datasetIdBytes, err := json.Marshal(payload.DatasetId)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "marshal dataset_id failed", "error": err})
			return
		}
		datasetIdJSON = string(datasetIdBytes)
	}

	result, err := difyClient.UpdateAppsInfo(payload.ID, payload.Name, payload.Description, payload.DatasetId, payload.Model, payload.PrePrompt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to Update Apps: %v", err)})
		return
	}

	var apps models.Apps
	resultDb := ac.DB.Where("id = ?", payload.ID).First(&apps)
	if resultDb.Error != nil {
		if resultDb.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "no record found", "error": resultDb.Error})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "database error", "error": resultDb.Error})
		}
		return
	}

	apps.Name = result.Name
	apps.Description = result.Description
	apps.PrePrompt = payload.PrePrompt
	apps.Model = payload.Model
	apps.Type = payload.Type
	apps.RoleIndustry = payload.RoleIndustry
	apps.RoleSettings = payload.RoleSettings
	apps.DatasetId = datasetIdJSON
	apps.UpdatedAt = time.Now()

	resultDb = ac.DB.Save(&apps)
	if resultDb.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "update apps failed", "error": resultDb.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": apps})
}

func (ac *AppsController) GetAppsForId(ctx *gin.Context) {
	id := ctx.Param("id")
	difyClient := InitDifyClient()

	result, err := difyClient.GetAppsHandler(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to Get Apps: %v", err)})
		return
	}
	appsRes := AppsForIdResponse{}

	appsRes.SiteStruct = result.Site
	appsRes.Datasets = result.ModelConfig.DatasetConfigs.Datasets.Datasets

	apps := models.Apps{}
	resultDb := ac.DB.Where("id = ?", id).First(&apps)
	if resultDb.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "find apps failed", "error": resultDb.Error})
		return
	}

	// Convert DatasetId to []string for response
	var datasetIDs []string
	if apps.DatasetId != "" {
		err := json.Unmarshal([]byte(apps.DatasetId), &datasetIDs)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "unmarshal dataset_id failed", "error": err})
			return
		}
	}

	appsRes.Apps = apps
	appsRes.DatasetId = datasetIDs // This will send the []string version of DatasetId in the response

	fmt.Println(appsRes.Apps.Description)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": appsRes})
}

func (ac *AppsController) GetAppsList(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	typeParam := ctx.Query("type")
	apps := []models.Apps{}

	var resultDb *gorm.DB
	if typeParam != "" {
		resultDb = ac.DB.Where("type = ? AND create_by = ?", typeParam, currentUser.ID).Order("updated_at DESC").Find(&apps)
	} else {
		resultDb = ac.DB.Where("create_by = ?", currentUser.ID).Order("updated_at DESC").Find(&apps)
	}

	if resultDb.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "find apps failed", "error": resultDb.Error})
		return
	}

	// Prepare response list with parsed dataset_id
	response := make([]AppsListResponse, len(apps))
	for i := range apps {
		response[i].Apps = apps[i]
		var datasetIDs []string
		if apps[i].DatasetId != "" {
			err := json.Unmarshal([]byte(apps[i].DatasetId), &datasetIDs)
			if err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "unmarshal dataset_id failed", "error": err})
				return
			}
		}
		response[i].DatasetId = datasetIDs
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func (ac *AppsController) UpdateIsPublished(ctx *gin.Context) {
	var payload struct {
		AppID       string `json:"app_id"`
		IsPublished bool   `json:"is_published"`
	}

	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询应用程序
	app := models.Apps{}
	resultDb := ac.DB.Where("id = ?", payload.AppID).First(&app)
	if resultDb.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "find app failed", "error": resultDb.Error})
		return
	}

	// 更新IsPublished字段
	resultDb = ac.DB.Model(&app).Update("is_published", payload.IsPublished)
	if resultDb.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "update IsPublished failed", "error": resultDb.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "IsPublished updated successfully"})
}

func (ac *AppsController) DeleteApp(ctx *gin.Context) {
	appID := ctx.Param("id")

	// 查询应用程序
	app := models.Apps{}
	resultDb := ac.DB.Where("id = ?", appID).First(&app)
	if resultDb.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "app not found", "error": resultDb.Error})
		return
	}

	difyClient := InitDifyClient()

	err := difyClient.DeleteAppsHandler(appID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to Delete Apps: %v", err)})
		return
	}

	// 删除应用程序
	resultDb = ac.DB.Delete(&app)
	if resultDb.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "delete app failed", "error": resultDb.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "app deleted successfully"})
}
