package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"hajime/golangp/common/utils"
	"net/http"
)

type ReferralCodeController struct {
	DB *gorm.DB
	cs *CreditSystem
}

func NewReferralCodeController(DB *gorm.DB, creditSystem *CreditSystem) ReferralCodeController {
	return ReferralCodeController{DB, creditSystem}
}

func (rc *ReferralCodeController) AddReferralCode(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	referralCodeModel := &models.ReferralCode{}
	referralCode, err := referralCodeModel.CreateReferralCode(rc.DB, currentUser.ID.String())

	if err != nil {
		logging.Warning("Failed to create referral code: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logging.Info("Created referral code: " + referralCode.Code + " for owner: " + currentUser.ID.String())
	ctx.JSON(http.StatusOK, gin.H{"referralCode": models.ValidateCode(referralCode.Code)})
	return
}

func (rc *ReferralCodeController) GetReferralCodeViaOwner(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	owner := currentUser.ID.String()

	referralCodes, err := models.GetReferralCodeViaOwner(rc.DB, owner)
	if err != nil {
		logging.Warning("Failed to get referral codes for owner: " + owner + " " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, referralCodes)
}

func (rc *ReferralCodeController) InviteUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	invitePayload := &models.InviteUserPayload{}

	if err := ctx.ShouldBindJSON(&invitePayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encodeCode := utils.Encode(invitePayload.Code)
	referralCode, err := models.GetReferralCode(rc.DB, encodeCode)

	if err != nil {
		logging.Warning("Failed to get referral code: " + encodeCode + " " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = referralCode.UpdateUsageCount(rc.DB)

	if err != nil {
		logging.Warning("Failed to update usage count for code: " + encodeCode + " " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.UpdateUserFrom(currentUser.ID.String(), referralCode.Code)

	if err != nil {
		logging.Warning("Failed to update user from: " + currentUser.ID.String() + " " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "referralCode": invitePayload.Code})
}

func (rc *ReferralCodeController) GetInvitedUserInfo(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Call the function to get invited users
	invitedUsersMap, err := models.GetInvitedUsersByReferralCode(rc.DB, currentUser.ID.String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invited users"})
		return
	}

	ctx.JSON(http.StatusOK, invitedUsersMap)
}
