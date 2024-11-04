package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"hajime/golangp/hajime_center/initializers"
	"hajime/golangp/hajime_center/models"
	"hajime/golangp/hajime_center/utils"
)

func DeserializeApp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var appToken string
		authorizationHeader := ctx.Request.Header.Get("App-Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			appToken = fields[1]
		}

		if appToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "App authorization token is missing"})
			return
		}

		config, _ := initializers.LoadEnv(".")
		sub, err := utils.ValidateToken(appToken, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var app models.Apps
		result := initializers.DB.First(&app, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the app belonging to this token no longer exists"})
			return
		}

		ctx.Set("currentApp", app)
		ctx.Next()
	}
}
