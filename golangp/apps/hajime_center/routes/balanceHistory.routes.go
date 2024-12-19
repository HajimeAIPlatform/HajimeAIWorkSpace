package routes

import (
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/middleware"

	"github.com/gin-gonic/gin"
)

type BalanceHistoryRouteController struct {
	balanceHistoryController controllers.BalanceHistoryController
}

func NewBalanceHistoryRouteController(balanceHistoryController controllers.BalanceHistoryController) BalanceHistoryRouteController {
	return BalanceHistoryRouteController{balanceHistoryController}
}

func (rc *BalanceHistoryRouteController) BalanceHistoryRoute(rg *gin.RouterGroup) {
	router := rg.Group("balance")
	router.GET("/get_balance_history", middleware.DeserializeUser(), rc.balanceHistoryController.GetBalanceHistoriesByUserID)
}
