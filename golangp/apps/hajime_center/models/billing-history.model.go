package models

import (
	"hajime/golangp/common/initializers"
	"time"

	"github.com/google/uuid"
)

type BillingHistory struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Operator          string    `gorm:"type:varchar(255);not null"`
	AccountEmail      string    `gorm:"not null"`
	Amount            int64     `gorm:"bigint;not null"`
	TransactionType   string    `gorm:"type:varchar(255)"`
	TransactionDetail string    `gorm:"type:varchar(255)"`
	TransactionTime   time.Time `gorm:"not null"`
}

func AddBillingHistory(email string, accountEmail, transactionType, transactionDetail string, amount int64) (*BillingHistory, error) {
	// Update the user's Address and Sign fields
	db := initializers.DB

	// Create a new BillingHistory instance
	billingHistory := &BillingHistory{
		Operator:          email,
		AccountEmail:      accountEmail,
		Amount:            amount,
		TransactionType:   transactionType,
		TransactionDetail: transactionDetail,
		TransactionTime:   time.Now(),
	}

	// Save the record to the database
	if err := db.Create(billingHistory).Error; err != nil {
		return nil, err
	}

	return billingHistory, nil
}
