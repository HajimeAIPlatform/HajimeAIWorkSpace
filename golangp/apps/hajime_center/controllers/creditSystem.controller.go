package controllers

import (
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
)

type CreditSystem struct {
	DB *gorm.DB
}

func NewCreditSystem(DB *gorm.DB) *CreditSystem {
	return &CreditSystem{DB}
}

func (cs *CreditSystem) GetBalanceByUserEmail(email string) (float64, error) {
	var user models.User
	err := cs.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return 0, err
	}
	logging.Info("Got balance of user", email, "as", user.Balance)
	return user.Balance, nil
}

func (cs *CreditSystem) UpdateBalanceByUserEmail(email string, amount float64) (float64, error) {
	var user models.User
	err := cs.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return 0, err
	}
	newBalance := user.Balance + amount
	err = cs.DB.Model(&user).Update("balance", newBalance).Error
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}
