package controllers

import (
	"hajime/golangp/apps/hajime_center/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BalanceHistoryController struct {
	DB *gorm.DB
}

func NewBalanceHistoryController(DB *gorm.DB) BalanceHistoryController {
	return BalanceHistoryController{DB}
}

func (bh *BalanceHistoryController) GetBalanceHistoriesByUserID(ctx *gin.Context) {
	// 获取当前用户
	currentUser, ok := ctx.MustGet("currentUser").(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "30")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid limit number"})
		return
	}

	// 调用模型方法
	result, err := models.GetBalanceHistoriesByUserID(currentUser.ID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error retrieving balance histories"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"data":     result.Data,
		"has_more": result.HasMore,
		"limit":    result.Limit,
		"page":     result.Page,
		"total":    result.Total,
	})
}
