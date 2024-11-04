package routes

import (
	"github.com/gin-gonic/gin"
	"hajime/golangp/apps/hajime_center/controllers"
	"hajime/golangp/apps/hajime_center/middleware"
)

type ConversationsRouteController struct {
	ConversationsController controllers.ConversationsController
}

func NewConversationsRouteController(conversationsController controllers.ConversationsController) ConversationsRouteController {
	return ConversationsRouteController{conversationsController}
}

func (dc *ConversationsRouteController) ConversationsRoute(rg *gin.RouterGroup) {
	router := rg.Group("conversations")
	router.Use(middleware.DeserializeUser())

	router.GET("/access_token", dc.ConversationsController.GetConversationsAccessToken)
	router.GET("/conversations", dc.ConversationsController.GetConversationsForApps)        //get all conversations
	router.DELETE("/:conversation_id", dc.ConversationsController.DeleteConversations)      //get all conversations
	router.POST("/:conversation_id/rename", dc.ConversationsController.RenameConversations) //get all conversations

	router.GET("/messages", dc.ConversationsController.GetConversationsMessages) //get conversation by id

	chatRouter := rg.Group("/chat")
	chatRouter.Use(middleware.DeserializeUser())
	chatRouter.POST("/chat-messages", dc.ConversationsController.ChatMessages)
	chatRouter.POST("/chat-messages/:task_id/stop", dc.ConversationsController.ChatMessagesStop)
}
