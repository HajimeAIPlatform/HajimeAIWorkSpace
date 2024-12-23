package models

import (
	"hajime/golangp/apps/hajime_center/initializers"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BalanceHistory struct {
	ID                      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID                  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ChangeBalance           float64   `gorm:"not null" json:"change_balance"`
	ChangeType              string    `gorm:"type:varchar(100);not null" json:"change_type"`
	Type                    string    `gorm:"type:varchar(100);not null" json:"type"`
	BalanceBefore           float64   `gorm:"not null" json:"balance_before"`
	BalanceAfter            float64   `gorm:"not null" json:"balance_after"`
	InviteBonusRateAddition float64   `gorm:"not null" json:"invite_bonus_rate_addition"`
	CreatedAt               time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type PaginatedBalanceHistories struct {
	Data    []BalanceHistory `json:"data"`
	HasMore bool             `json:"has_more"`
	Limit   int              `json:"limit"`
	Page    int              `json:"page"`
	Total   int64            `json:"total"`
}

func AddBalanceHistory(db *gorm.DB, userID uuid.UUID, changeBalance float64, changeType string, operatorType string, balanceBefore float64, balanceAfter float64, inviteBonusRateAddition float64) error {
	// 检查 changeType 是否有效

	// 创建新的 BalanceHistory 记录
	history := BalanceHistory{
		ID:                      uuid.New(),
		UserID:                  userID,
		ChangeBalance:           changeBalance,
		ChangeType:              changeType,
		Type:                    operatorType,
		BalanceBefore:           balanceBefore,
		BalanceAfter:            balanceAfter,
		InviteBonusRateAddition: inviteBonusRateAddition,
		CreatedAt:               time.Now(),
	}

	// 将记录插入数据库
	if err := db.Create(&history).Error; err != nil {
		return err
	}

	return nil
}

func GetBalanceHistoriesByUserID(userID uuid.UUID, page, pageSize int) (PaginatedBalanceHistories, error) {
	db := initializers.DB
	var histories []BalanceHistory
	var totalCount int64

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query to count total records for the user
	if err := db.Model(&BalanceHistory{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return PaginatedBalanceHistories{}, err
	}

	// Query to get the paginated records
	if err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&histories).Error; err != nil {
		return PaginatedBalanceHistories{}, err
	}

	// Determine if there is a next page
	hasMore := int64(offset+pageSize) < totalCount

	return PaginatedBalanceHistories{
		Data:    histories,
		HasMore: hasMore,
		Limit:   pageSize,
		Page:    page,
		Total:   totalCount,
	}, nil
}
