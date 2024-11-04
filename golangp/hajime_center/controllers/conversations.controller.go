package controllers

import (
	"hajime/golangp/hajime_center/dify"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ConversationsController struct {
	DB *gorm.DB
}

type RenameConversationsPayload struct {
	Name string `json:"name"`
}

func NewConversationsController(DB *gorm.DB) ConversationsController {
	return ConversationsController{DB}
}

func (c *ConversationsController) GetConversationsAccessToken(ctx *gin.Context) {
	accessToken := ctx.Query("access_token")

	difyClient := InitDifyClient()

	result, err := difyClient.GetAppsAccessToken(accessToken)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": result.AccessToken})
}

func GetAppAuthorization(ctx *gin.Context) string {
	return ctx.GetHeader("APP-Authorization")
}

func (c *ConversationsController) GetConversationsForApps(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	appAuthorization := GetAppAuthorization(ctx)

	if appAuthorization == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid limit parameter"})
		return
	}

	difyClient := InitDifyClient()

	result, err := difyClient.GetConversations(appAuthorization, limit)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result.Data, "hasMore": result.HasMore, "limit": result.Limit})
}

func (c *ConversationsController) GetConversationsMessages(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	lastId := ctx.Query("last_id")
	conversationId := ctx.Query("conversation_id")

	appAuthorization := GetAppAuthorization(ctx)

	if appAuthorization == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid limit parameter"})
		return
	}

	difyClient := InitDifyClient()

	result, err := difyClient.GetConversationsMessages(appAuthorization, lastId, conversationId, limit)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result.Data, "hasMore": result.HasMore, "limit": result.Limit})
}

func (c *ConversationsController) ChatMessages(ctx *gin.Context) {
	payload := dify.ChatMessagesPayload{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	appAuthorization := GetAppAuthorization(ctx)

	if appAuthorization == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	difyClient := InitDifyClient()

	if payload.ResponseMode == "blocking" {
		result, err := difyClient.ChatMessages(payload.Query, payload.Inputs, payload.ConversationID, payload.Files, appAuthorization)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
	} else {
		// 设置流式响应头
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")

		dataChan, err := difyClient.ChatMessagesStreaming(payload.Query, payload.Inputs, payload.ConversationID, payload.Files, appAuthorization)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		// 开始流式传输数据
		for msg := range dataChan {
			_, err := ctx.Writer.Write([]byte(msg + "\n\n"))
			if err != nil {
				break
			}

			// 刷新缓冲区，确保数据被发送到客户端
			ctx.Writer.Flush()
		}
	}
}

func (c *ConversationsController) ChatMessagesStop(ctx *gin.Context) {
	task_id := ctx.Param("task_id") // 从URL参数获取 task_id
	if task_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing task_id parameter"})
		return
	}

	AppAuthorization := GetAppAuthorization(ctx)
	if AppAuthorization == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	difyClient := InitDifyClient()
	_, err := difyClient.ChatMessagesStop(task_id, AppAuthorization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (c *ConversationsController) DeleteConversations(ctx *gin.Context) {
	conversationId := ctx.Param("conversation_id") // 从URL参数获取 conversation_id

	appAuthorization := GetAppAuthorization(ctx)
	if appAuthorization == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	difyClient := InitDifyClient()
	err := difyClient.DeleteConversations(conversationId, appAuthorization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (c *ConversationsController) RenameConversations(ctx *gin.Context) {
	conversation_id := ctx.Param("conversation_id") // 从URL参数获取 conversation_id
	if conversation_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing conversation_id parameter"})
		return
	}

	payload := RenameConversationsPayload{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	name := payload.Name

	appAuthorization := GetAppAuthorization(ctx)
	if appAuthorization == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing APP-Authorization header"})
		return
	}

	difyClient := InitDifyClient()
	result, err := difyClient.RenameConversations(conversation_id, name, appAuthorization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}
