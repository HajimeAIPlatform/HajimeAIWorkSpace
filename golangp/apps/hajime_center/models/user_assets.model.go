package models

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/initializers"
	"time"
)

type UserAsset struct {
	ID        uint            `gorm:"primaryKey;autoIncrement"`
	UID       string          `gorm:"type:varchar(255);not null"`
	Mainchain string          `gorm:"type:varchar(255);not null;default:'SOLANA'"`
	Token     string          `gorm:"type:varchar(255);not null;default:'SOL'"`
	Amount    decimal.Decimal `gorm:"type:decimal(20,8);not null;default:0.00000000"`
	Frozen    decimal.Decimal `gorm:"type:decimal(20,8);not null;default:0.00000000"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
}

type UserWithdraw struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UID       string `gorm:"type:varchar(255);not null"`
	Address   string `gorm:"type:varchar(255);not null"`
	Mainchain string `gorm:"type:varchar(255);not null"`
	Token     string `gorm:"type:varchar(255);not null"`
	Amount    decimal.Decimal
	Desc      string `gorm:"type:varchar(255)"`
}

func (u *UserAsset) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *UserAsset) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

func (u *UserAsset) InitToken(uid, mainchain, token string) (*UserAsset, error) {
	db := initializers.DB
	var asset UserAsset
	err := db.Where("uid = ? AND mainchain = ? AND token = ?", uid, mainchain, token).First(&asset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		asset = UserAsset{
			UID:       uid,
			Mainchain: mainchain,
			Token:     token,
			Amount:    decimal.NewFromFloat(0),
			Frozen:    decimal.NewFromFloat(0),
		}

		if err := db.Create(&asset).Error; err != nil {
			return nil, err
		}
		return &asset, nil
	} else if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (u *UserAsset) GetUserAssetList(uid string) ([]UserAsset, error) {
	db := initializers.DB
	var assets []UserAsset
	err := db.Where("uid = ?", uid).Find(&assets).Error
	return assets, err
}

func (u *UserAsset) GetUserAssetMap(uid string) (map[string]decimal.Decimal, error) {
	assets, err := u.GetUserAssetList(uid)
	if err != nil {
		return nil, err
	}

	out := make(map[string]decimal.Decimal)
	for _, asset := range assets {
		out[asset.Token] = asset.Amount
	}
	return out, nil
}

func (u *UserAsset) GetUserFilterAssetMap(where map[string]interface{}) (map[string]decimal.Decimal, error) {
	db := initializers.DB
	var assets []UserAsset
	err := db.Where(where).Find(&assets).Error
	if err != nil {
		return nil, err
	}

	out := make(map[string]decimal.Decimal)
	for _, asset := range assets {
		out[asset.Token] = asset.Amount
	}
	return out, nil
}

func (u *UserAsset) Incr(uid, mainchain, token string, amount decimal.Decimal, opType, desc string) (interface{}, error) {
	db := initializers.DB
	var asset UserAsset
	err := db.Where("uid = ? AND mainchain = ? AND token = ?", uid, mainchain, token).First(&asset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		asset = UserAsset{
			UID:       uid,
			Mainchain: mainchain,
			Token:     token,
			Amount:    decimal.NewFromFloat(0),
			Frozen:    decimal.NewFromFloat(0),
		}
		if err := db.Create(&asset).Error; err != nil {
			return nil, err
		}
	}

	if amount.LessThan(decimal.NewFromFloat(0)) {
		if asset.Amount.LessThan(amount.Neg()) {
			return nil, fmt.Errorf("insufficient balance")
		}
	}

	err = db.Model(&asset).Where("uid = ? AND mainchain = ? AND token = ?", uid, mainchain, token).
		Update("amount", gorm.Expr("amount + ?", amount)).Error
	if err != nil {
		return nil, err
	}

	// Log the transaction if necessary

	return asset, nil
}

func (uw *UserWithdraw) AddWithdraw(db *gorm.DB) error {
	return db.Create(uw).Error
}

func (u *UserAsset) Withdraw(uid, address, token string, num float64, opType, desc, mainchain string) (interface{}, error) {
	db := initializers.DB
	ret, err := u.Incr(uid, mainchain, token, decimal.NewFromFloat(-num), opType, desc)
	if err != nil {
		return nil, err
	}

	// Check if the increment operation was successful
	if ret != nil {
		withdraw := UserWithdraw{
			UID:       uid,
			Address:   address,
			Mainchain: mainchain,
			Token:     token,
			Amount:    decimal.NewFromFloat(num),
			Desc:      "user withdraw",
		}

		if err := withdraw.AddWithdraw(db); err != nil {
			return nil, fmt.Errorf("failed to add withdrawal record: %v", err)
		}

		return num, nil
	} else {
		return nil, fmt.Errorf("insufficient balance")
	}
}
