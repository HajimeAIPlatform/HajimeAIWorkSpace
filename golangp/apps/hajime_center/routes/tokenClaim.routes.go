package routes

import (
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/middleware"

	"github.com/gin-gonic/gin"
)

type TokenClaimRouteController struct {
	tokenClaimController controllers.TokenClaimController // 使用值类型
}

func NewTokenClaimRouteController(csvFilePath string) TokenClaimRouteController {
	// 返回值类型的 TokenClaimRouteController 实例
	tokenClaimController := controllers.NewTokenClaimController(csvFilePath)
	return TokenClaimRouteController{tokenClaimController: tokenClaimController}
}

func (rc *TokenClaimRouteController) TokenClaimRoute(rg *gin.RouterGroup) {
	router := rg.Group("token_claim")
	// router.GET("/get_info/:address", rc.tokenClaimController.GetSolanaAddressInfo)
	router.GET("/get_info/:address", middleware.DeserializeUser(), rc.tokenClaimController.GetSolanaAddressInfo)
}
