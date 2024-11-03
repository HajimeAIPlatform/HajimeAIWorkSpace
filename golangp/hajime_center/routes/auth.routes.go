package routes

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/controllers"
	"HajimeAIWorkSpace/common/apps/hajime_center/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("user")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", rc.authController.LogoutUser)
	router.GET("/verifyemail/:verificationCode", rc.authController.VerifyEmail)
	router.POST("/forgotpassword", rc.authController.ForgotPassword)
	router.PATCH("/resetpassword/:resetToken", rc.authController.ResetPassword)
	router.POST("/password", middleware.DeserializeUser(), rc.authController.PasswordChange)

	// admin manager routes
	router.POST("/add", middleware.DeserializeUser(), rc.authController.AddUser)
	router.DELETE("/:userId", middleware.DeserializeUser(), rc.authController.DeleteUser)
	router.GET("/users", middleware.DeserializeUser(), rc.authController.GetAllUsers)
	router.PUT("/:userId", middleware.DeserializeUser(), rc.authController.UpdateUser)

	// Only apply middleware.DeserializeUser() to specific routes
	router.GET("/getme", middleware.DeserializeUser(), rc.authController.GetMe)
}
