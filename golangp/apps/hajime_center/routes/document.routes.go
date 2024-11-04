package routes

import (
	"github.com/gin-gonic/gin"
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/middleware"
)

type DocumentRouteController struct {
	documentController controllers.DocumentController
}

func NewDocumentRouteController(documentController controllers.DocumentController) DocumentRouteController {
	return DocumentRouteController{documentController}
}

func (dc *DocumentRouteController) DocumentRoute(rg *gin.RouterGroup) {

	router := rg.Group("documents")
	router.Use(middleware.DeserializeUser())

	router.GET("/accessToken/test", dc.documentController.GetDidyAccessToken)

	fileRouter := rg.Group("file")
	fileRouter.Use(middleware.DeserializeUser())
	fileRouter.POST("/upload", dc.documentController.HandleFileUploadWithDatasets)
	fileRouter.POST("/upload/chat", dc.documentController.HandleFileUploadForChat)
	fileRouter.GET("/supportType", dc.documentController.GetFilesSupportType)

	fileRouter.GET("/datasets/:dataset_id", dc.documentController.GetDatasetFileList)
	fileRouter.GET("/datasets/:dataset_id/indexing_status", dc.documentController.GetDatasetIndexingStatus)
	fileRouter.GET("/datasets/:dataset_id/:document_id/indexing_status", dc.documentController.GetDocumentIndexingStatus)

	fileRouter.DELETE("/datasets/:dataset_id/:document_id", dc.documentController.DeleteDocumentForDataset)
	fileRouter.POST("/datasets/:dataset_id/:document_id/rename", dc.documentController.RenameDocumentForDataset)

	rg.GET("/file/:file_id/preview", dc.documentController.PreviewFile)
	rg.GET("/file/preview/:file_id", dc.documentController.PreviewFileByURL)
}
