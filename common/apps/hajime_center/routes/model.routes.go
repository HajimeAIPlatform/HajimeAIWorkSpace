package routes

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/controllers"
	"HajimeAIWorkSpace/common/apps/hajime_center/middleware"
	"github.com/gin-gonic/gin"
)

type ModelRouteController struct {
	ModelController controllers.ModelController
}

func NewModelRouteController(modelController controllers.ModelController) ModelRouteController {
	return ModelRouteController{modelController}
}

func (dc *ModelRouteController) ModelRoute(rg *gin.RouterGroup) {
	router := rg.Group("model")
	router.Use(middleware.DeserializeUser())

	router.GET("/default-model", dc.ModelController.GetModelsDefault)
	router.GET("/model", dc.ModelController.GetAllModels)
}
