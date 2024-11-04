package controllers

import (
	"hajime/golangp/hajime_center/dify"
	"hajime/golangp/hajime_center/initializers"
	"hajime/golangp/hajime_center/logger"
	"hajime/golangp/hajime_center/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

type DocumentController struct {
	DB     *gorm.DB
	DBDify *gorm.DB
}

func NewDocumentController(DB *gorm.DB, DBDify *gorm.DB) DocumentController {
	return DocumentController{DB, DBDify}
}

func InitDifyClient() *dify.DifyClient {
	client, err := dify.GetDifyClient()
	if err != nil {
		logger.Warning(err.Error())
		return nil // 返回 nil 以符合返回类型
	}

	_, err = client.GetUserToken()
	if err != nil {
		logger.Warning(err.Error())
		return nil // 返回 nil 以符合返回类型
	}
	return client
}

func (c *DocumentController) GetDidyAccessToken(ctx *gin.Context) {
	client, err := dify.GetDifyClient()

	if err != nil {
		logger.Warning(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to get client"})
		return
	}

	accessToken, err := client.GetUserToken()
	if err != nil {
		logger.Warning(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to get access token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

// 定义一个新的结构体来包含所需的所有信息
type CustomResponse struct {
	dify.FileUploadResponse
	dify.InitDatasetsResponse
}

type RenameDocumentInput struct {
	Name string `json:"name"`
}

func (c *DocumentController) HandleFileUploadWithDatasets(ctx *gin.Context) {
	// 获取文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

	difyClient := InitDifyClient()

	// 上传文件
	result, err := difyClient.DatasetsFileUpload(file, header.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// 创建一个新的 CustomResponse 变量
	customResponse := CustomResponse{
		FileUploadResponse: result,
	}

	// 返回上传结果
	fileId := result.ID
	datasets_id := ctx.PostForm("datasets_id")

	var resultDataset dify.InitDatasetsResponse
	var errDataset error
	if datasets_id != "" {
		resultDataset, errDataset = difyClient.UploadDatasetsByUploadFile([]string{fileId}, datasets_id)
	} else {
		resultDataset, errDataset = difyClient.InitDatasetsByUploadFile([]string{fileId})
	}

	if errDataset != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to init datasets: %v", err)})
		return
	}
	customResponse.Dataset = resultDataset.Dataset
	customResponse.Documents = resultDataset.Documents
	customResponse.Batch = resultDataset.Batch

	// 返回上传结果
	ctx.JSON(http.StatusOK, customResponse)
}

func (c *DocumentController) GetFilesSupportType(ctx *gin.Context) {
	difyClient := InitDifyClient()

	// 上传文件
	result, err := difyClient.GetSupportTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// 返回上传结果
	ctx.JSON(http.StatusOK, result)
}

func (c *DocumentController) HandleFileUploadForChat(ctx *gin.Context) {
	// 获取文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

	difyClient := InitDifyClient()

	// 上传文件
	result, err := difyClient.DatasetsFileUploadChat(file, header.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *DocumentController) GetDatasetFileList(ctx *gin.Context) {
	dataset_id := ctx.Param("dataset_id")

	if dataset_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid dataset_id parameter"})
		return
	}

	limitStr := ctx.Query("limit")
	pageStr := ctx.Query("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid limit parameter"})
		return
	}

	difyClient := InitDifyClient()
	result, err := difyClient.GetDatasetsFileList(dataset_id, limit, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get dataset file list: %v", err)})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DocumentController) GetDatasetIndexingStatus(ctx *gin.Context) {
	dataset_id := ctx.Param("dataset_id")

	if dataset_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid dataset_id parameter"})
		return
	}

	difyClient := InitDifyClient()
	result, err := difyClient.InitDatasetsIndexingStatus(dataset_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get dataset file indexing status: %v", err)})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DocumentController) GetDocumentIndexingStatus(ctx *gin.Context) {
	dataset_id := ctx.Param("dataset_id")
	document_id := ctx.Param("document_id")

	if dataset_id == "" || document_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid dataset_id parameter"})
		return
	}

	difyClient := InitDifyClient()
	result, err := difyClient.DocumentIndexingStatus(dataset_id, document_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get dataset file indexing status: %v", err)})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DocumentController) DeleteDocumentForDataset(ctx *gin.Context) {
	dataset_id := ctx.Param("dataset_id")
	document_id := ctx.Param("document_id")
	difyClient := InitDifyClient()
	_, err := difyClient.DeleteDocumentForDatasets(dataset_id, document_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete dataset file: %v", err)})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (c *DocumentController) RenameDocumentForDataset(ctx *gin.Context) {
	payload := RenameDocumentInput{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dataset_id := ctx.Param("dataset_id")
	document_id := ctx.Param("document_id")
	difyClient := InitDifyClient()
	result, err := difyClient.RenameDocumentForDatasets(dataset_id, document_id, payload.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete dataset file: %v", err)})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *DocumentController) PreviewFile(ctx *gin.Context) {
	fileID := ctx.Param("file_id")

	difyClient := InitDifyClient()

	body, contentType, err := difyClient.PreviewFile(fileID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Header("Content-Type", contentType)
	ctx.Status(http.StatusOK)
	ctx.Writer.Write(body)
}

func (c *DocumentController) PreviewFileByURL(ctx *gin.Context) {
	fileID := ctx.Param("file_id")

	fileResult, err := models.QueryStorageByID(c.DBDify, fileID)

	fmt.Printf("fileResult: %v\n", fileResult.Key)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	config, _ := initializers.LoadEnv(".")

	filePath := config.DifyConsoleStoragePath + "/" + fileResult.Key

	fmt.Printf("filePath: %v\n", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.String(http.StatusNotFound, "File not found")
		return
	}

	// // 只允许特定的文件类型
	// allowedExtensions := map[string]bool{
	// 	"png":  true,
	// 	"jpg":  true,
	// 	"jpeg": true,
	// 	"pdf":  true,
	// }

	// if !allowedExtensions[extension] {
	// 	ctx.String(http.StatusForbidden, "File type not allowed")
	// 	return
	// }

	ctx.File(filePath)
}
