package routes

import (
	"HajimeAIWorkSpace/common/apps/hajime_center/controllers"
	"HajimeAIWorkSpace/common/apps/hajime_center/middleware"
	"github.com/gin-gonic/gin"
)

type ChatRouteController struct {
	chatController controllers.ChatController
}

func NewChatRouteController(chatController controllers.ChatController) ChatRouteController {
	return ChatRouteController{chatController}
}

func (crc *ChatRouteController) ChatRoute(rg *gin.RouterGroup) {

	chat := rg.Group("chat").Use(middleware.DeserializeUser())
	{
		chat.POST("/completion_with_model_info", crc.chatController.CompletionWithModelInfo)
	}
}
