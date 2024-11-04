package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ModelController struct {
	DB *gorm.DB
}

func NewModelController(DB *gorm.DB) ModelController {
	return ModelController{DB}
}

func (c *ModelController) GetModelsDefault(ctx *gin.Context) {
	difyClient := InitDifyClient()

	result, err := difyClient.GetCurrentWorkspaceLLMDefaultModel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	fmt.Println(result)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result.Data})
}

func (c *ModelController) GetAllModels(ctx *gin.Context) {
	difyClient := InitDifyClient()

	result, err := difyClient.GetCurrentWorkspaceLLMModel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	fmt.Println(result)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result.Data})
}
