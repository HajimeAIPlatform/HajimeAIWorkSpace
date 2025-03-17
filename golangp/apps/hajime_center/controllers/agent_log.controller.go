package controllers

import (
	"hajime/golangp/apps/hajime_center/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AgentLogController handles log-related HTTP requests
type AgentLogController struct {
	service *services.AgentLogService
}

// NewAgentLogController creates a new instance of AgentLogController
func NewAgentLogController(service *services.AgentLogService) *AgentLogController {
	return &AgentLogController{service: service}
}

// GetAgentLogs handles GET /agent_logs request
func (c *AgentLogController) GetAgentLogs(ctx *gin.Context) {
	// Get query parameters
	logType := ctx.Query("type")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// Fetch logs from service
	logs, total, err := c.service.GetLogs(logType, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Prepare response
	ctx.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RegisterRoutes registers controller routes
func (c *AgentLogController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/agent_logs", c.GetAgentLogs)
}
