package routes

import (
	"github.com/gin-gonic/gin"
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/middleware"
)

type ReferralCodeRouteController struct {
	referralCodeController controllers.ReferralCodeController
}

func NewReferralCodeRouteController(referralCodeController controllers.ReferralCodeController) ReferralCodeRouteController {
	return ReferralCodeRouteController{referralCodeController}
}

func (rc *ReferralCodeRouteController) ReferralCodeRoute(rg *gin.RouterGroup) {
	router := rg.Group("referral_code")
	router.POST("/add_code", middleware.DeserializeUser(), rc.referralCodeController.AddReferralCode)
	router.GET("/get_user_referral_code", middleware.DeserializeUser(), rc.referralCodeController.GetReferralCodeViaOwner)
	router.POST("/invite_user", middleware.DeserializeUser(), rc.referralCodeController.InviteUser)
	router.GET("/invited_user_info", middleware.DeserializeUser(), rc.referralCodeController.GetInvitedUserInfo)
}
