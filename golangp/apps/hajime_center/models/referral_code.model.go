package models

import (
	"fmt"
	"hajime/golangp/common/utils"
	"time"

	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type ReferralCode struct {
	Code       string `gorm:"type:varchar(255);primaryKey" json:"code"`
	OriginCode string `gorm:"type:varchar(255)" json:"origin_code"`
	Owner      string `gorm:"type:varchar(255)" json:"owner"`
	UsageCount int    `gorm:"default:0" json:"usage_count"`
	CreatedAt  int64  `json:"created_at"`
}

type InviteUserInfo struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FromCode string    `json:"from_code,omitempty"`
}

type InviteUserPayload struct {
	Code string `json:"code" binding:"required"`
}

type InvitedUsersResponse struct {
	Users      []string `json:"users"`
	UsageCount int      `json:"usage_count"`
}

func (rc *ReferralCode) GenerateRandomCode() (referralCode string, originCode string) {
	code := randstr.String(10)

	referralCode = utils.Encode(code)
	return referralCode, code
}

func ValidateCode(code string) (referralCode string) {
	referralCode, _ = utils.Decode(code)
	return referralCode
}

// CreateReferralCode adds a new ReferralCode to the database
func (rc *ReferralCode) CreateReferralCode(db *gorm.DB, user User) (*ReferralCode, error) {
	if user.UsedCodeAmount >= user.UserMaxCodeAmount {
		return nil, fmt.Errorf("user has reached the maximum number of invite code")
	}
	code, originCode := rc.GenerateRandomCode()
	referralCode := &ReferralCode{
		Code:       code,
		OriginCode: originCode,
		Owner:      user.ID.String(),
		UsageCount: 0,
		CreatedAt:  time.Now().Unix(),
	}

	if err := db.Create(referralCode).Error; err != nil {
		return nil, err
	}
	err := UpdateUserCodeAmount(user.ID.String(), 1)
	if err != nil {
		return nil, err
	}

	return referralCode, nil
}

func (rc *ReferralCode) UpdateUsageCount(db *gorm.DB) error {
	return db.Model(&ReferralCode{}).Where("code = ?", rc.Code).Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func GetReferralCode(db *gorm.DB, code string) (*ReferralCode, error) {
	var referralCode ReferralCode
	if err := db.Where("code = ?", code).First(&referralCode).Error; err != nil {
		return nil, err
	}
	return &referralCode, nil
}

func GetReferralCodeViaOwner(db *gorm.DB, owner string) ([]ReferralCode, error) {
	var referralCodes []ReferralCode
	if err := db.Where("owner = ?", owner).Find(&referralCodes).Error; err != nil {
		return nil, err
	}
	for i := range referralCodes {
		referralCodes[i].Code = ValidateCode(referralCodes[i].Code)
	}
	return referralCodes, nil
}

func GetInvitedUsersByOwnerId(db *gorm.DB, userId string) (map[string]InvitedUsersResponse, error) {
	// Step 1: Retrieve all referral codes owned by the user
	var referralCodes []ReferralCode
	if err := db.Where("owner = ?", userId).Find(&referralCodes).Error; err != nil {
		return nil, err
	}

	// Step 2: Initialize a map to hold the results
	invitedUsersMap := make(map[string]InvitedUsersResponse)

	// Step 3: Retrieve all users with a from_code
	var allUsers []struct {
		Email    string
		FromCode string
	}
	if err := db.Model(&User{}).Select("email, from_code").Find(&allUsers).Error; err != nil {
		return nil, err
	}

	// Step 4: Filter users by referral codes and calculate usage count
	for _, user := range allUsers {
		for _, rc := range referralCodes {
			code := ValidateCode(rc.Code)
			if user.FromCode == rc.Code {
				response := invitedUsersMap[code]
				response.Users = append(response.Users, user.Email)
				response.UsageCount = rc.UsageCount
				invitedUsersMap[code] = response
			}
		}
	}

	return invitedUsersMap, nil
}
