package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BalanceHistory struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ChangeBalance float64   `gorm:"not null" json:"change_balance"`
	ChangeType    string    `gorm:"type:varchar(100);not null" json:"change_type"`
	Type          string    `gorm:"type:varchar(100);not null" json:"type"`
	BalanceBefore float64   `gorm:"not null" json:"balance_before"`
	BalanceAfter  float64   `gorm:"not null" json:"balance_after"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func AddBalanceHistory(db *gorm.DB, userID uuid.UUID, changeBalance float64, changeType string, operatorType string, balanceBefore float64, balanceAfter float64) error {
	// 检查 changeType 是否有效

	// 创建新的 BalanceHistory 记录
	history := BalanceHistory{
		ID:            uuid.New(),
		UserID:        userID,
		ChangeBalance: changeBalance,
		ChangeType:    changeType,
		Type:          operatorType,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     time.Now(),
	}

	// 将记录插入数据库
	if err := db.Create(&history).Error; err != nil {
		return err
	}

	return nil
}
