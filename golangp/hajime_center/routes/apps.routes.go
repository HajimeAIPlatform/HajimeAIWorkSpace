package routes

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/controllers"
	"HajimeAIWorkSpace/common/apps/hajime_center/middleware"
	"github.com/gin-gonic/gin"
)

type AppsRouteController struct {
	AppsController controllers.AppsController
}

func NewAppsRouteController(appsController controllers.AppsController) AppsRouteController {
	return AppsRouteController{appsController}
}

func (dc *AppsRouteController) AppsRoute(rg *gin.RouterGroup) {
	router := rg.Group("app")
	router.Use(middleware.DeserializeUser())

	router.POST("/create", dc.AppsController.CreateApps)
	router.POST("/update", dc.AppsController.UpdateApps)
	router.GET("/getAll", dc.AppsController.GetAppsList)
	router.GET("/:id", dc.AppsController.GetAppsForId)
	router.DELETE("/:id", dc.AppsController.DeleteApp)
	router.POST("/publish", dc.AppsController.UpdateIsPublished)
}
