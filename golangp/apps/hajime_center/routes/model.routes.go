package routes

import (
	"github.com/gin-gonic/gin"
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/middleware"
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
