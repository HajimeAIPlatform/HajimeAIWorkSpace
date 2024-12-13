package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/common/logging"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name               string    `gorm:"type:varchar(255);not null"`
	Email              string    `gorm:"uniqueIndex;not null"`
	Password           string    `gorm:"not null"`
	Role               string    `gorm:"type:varchar(255);not null"`
	Provider           string    `gorm:"not null"`
	Photo              string    `gorm:"not null"`
	VerificationCode   string
	PasswordResetToken string
	PasswordResetAt    time.Time
	Verified           bool       `gorm:"not null"`
	Balance            float64    `gorm:"not null;default:1000"`
	Address            string     `gorm:"type:varchar(255);default:''"`
	Sign               string     `gorm:"type:varchar(255);default:''"`
	Status             int32      `gorm:"not null;default:1"` // Corrected type
	Code               string     `gorm:"type:varchar(255);default:''"`
	Twitter            string     `gorm:"type:varchar(255);default:''"` // Twitter
	Telegram           string     `gorm:"type:varchar(255);default:''"` // Telegram
	Discord            string     `gorm:"type:varchar(255);default:''"` // Discord
	AppPublishAmount   int64      `gorm:"not null;default:0"`
	AppAmount          int64      `gorm:"not null;default:0"`
	CreatedAt          time.Time  `gorm:"not null"`
	UpdatedAt          time.Time  `gorm:"not null"`
	FromCode           string     `gorm:"type:varchar(255)"`
	UsedCodeAmount     int        `gorm:"not null;default:0"`
	UserMaxCodeAmount  int        `gorm:"not null;default:0"`
	LoginTime          *time.Time `gorm:"default:null"`
	AppUsage           string     `gorm:"type:jsonb;default:'{}'"`
	ConfigUsage        string     `gorm:"type:jsonb;default:'{}'"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	Photo           string `json:"photo,omitempty"`
	Role            string `json:"role"`
	FromCode        string `json:"fromCode,omitempty"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	Verified  bool      `json:"verified"`
	Balance   float64   `json:"balance,default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ForgotPasswordInput struct
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

// ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

type UpdateBalanceInput struct {
	Email  string `json:"email"  binding:"required"`
	Amount int64  `json:"amount"  binding:"required"`
}

type PasswordChangeInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,min=8"`
}

type UpdateUserInput struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
}

type SignLoginModel struct {
	WalletAddress string `json:"walletAddress"`
	Sign          string `json:"sign"`
	Code          string `json:"code,omitempty"`
	Msg           string `json:"msg,omitempty"`
}

type LoginResponse struct {
	WalletAddress string `json:"walletAddress"`
	UID           string `json:"uid,omitempty"`
	Sign          string `json:"sign,omitempty"`
}

func (u *User) GetUserByAddress(address string) (*User, error) {
	db := initializers.DB
	var user User
	result := db.Where("address = ?", address).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *User) UpdateAppUsage(appID string) (bool, error) {
	db := initializers.DB
	exceedsLimit := false
	var usageData map[string]int

	if u.AppUsage != "" {
		// Parse JSON string into map
		if err := json.Unmarshal([]byte(u.AppUsage), &usageData); err != nil {
			return false, err
		}
	} else {
		usageData = make(map[string]int)
	}

	// Check if the appID exists in the map
	if count, ok := usageData[appID]; ok {
		usageData[appID] = count + 1
	} else {
		usageData[appID] = 1
		err := u.UpdateBalance(constants.UseBotAgentWorkflowPoints, "UseBotAgentWorkflowPoints") // Replace with your actual points
		if err != nil {
			return false, err
		}
	}

	// Check if the value exceeds 3
	if usageData[appID] > 3 {
		exceedsLimit = true

		isCreditsEnough := u.PreCheckBalance()

		if !isCreditsEnough {
			return false, errors.New("score not enough, you currently have " + fmt.Sprint(u.Balance) + " score")
		}
	}

	// Convert map back to JSON
	updatedAppUsage, err := json.Marshal(usageData)
	if err != nil {
		return false, err
	}
	u.AppUsage = string(updatedAppUsage)

	// Save the user to the database
	result := db.Save(u)
	if result.Error != nil {
		return exceedsLimit, result.Error
	}

	return exceedsLimit, nil
}

func (u *User) UpdateConfigUsage(appID string, knowledge bool, variables bool) error {
	db := initializers.DB
	var configUsage map[string][]string

	if u.ConfigUsage != "" {
		// Parse JSON string into map
		if err := json.Unmarshal([]byte(u.ConfigUsage), &configUsage); err != nil {
			return err
		}
	} else {
		configUsage = make(map[string][]string)
	}

	//Initialize or update the appID entry
	if _, ok := configUsage[appID]; !ok {
		configUsage[appID] = []string{"", ""}
	}

	// Check and update "Knowledge"
	if knowledge && configUsage[appID][0] == "" {
		configUsage[appID][0] = "Knowledge"
		err := u.UpdateBalance(constants.UseKnowledgePoints, "UseKnowledgePoints")
		if err != nil {
			return err
		}
	}

	// Check and update "Variables"
	if variables && configUsage[appID][1] == "" {
		configUsage[appID][1] = "Variables"
		err := u.UpdateBalance(constants.UseVariablesPoints, "UseVariablesPoints")
		if err != nil {
			return err
		}
	}

	// Convert map back to JSON
	updatedConfigUsage, err := json.Marshal(configUsage)
	if err != nil {
		return err
	}
	u.ConfigUsage = string(updatedConfigUsage)

	// Save the user to the database
	result := db.Save(u)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *User) SaveUser(user *User) error {
	db := initializers.DB
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) GetUserByUID(uid string) (*User, error) {
	db := initializers.DB
	var user User
	result := db.Where("id = ?", uid).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, result.Error // Other errors
	}
	return &user, nil
}

func (u *User) UpdateAddressAndSign(address, sign string) error {
	// Update the user's Address and Sign fields
	db := initializers.DB
	result := db.Model(u).Updates(User{Address: address, Sign: sign})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) UpdateLoginTime(loginTime *time.Time) error {
	// Update the user's LoginTime field
	db := initializers.DB
	result := db.Model(u).Update("LoginTime", loginTime)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (u *User) UpdateBalance(balance float64, operatorType string) error {
	// Update the user's Balance field
	db := initializers.DB
	// 确定变动类型
	changeType := "add"
	if balance < 0 {
		changeType = "subtract"
	}
	newBalance := u.Balance + balance

	err := AddBalanceHistory(db, u.ID, balance, changeType, operatorType, u.Balance, newBalance)
	if err != nil {
		return err
	}

	result := db.Model(u).Update("Balance", newBalance)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) GenerateLoginResponse(walletAddress, uid string) *LoginResponse {
	return &LoginResponse{
		WalletAddress: walletAddress,
		UID:           uid,
	}
}

func (u *User) LoginWithSign(form SignLoginModel, UID uuid.UUID) (*LoginResponse, error) {
	uidStr := UID.String()

	user, err := u.GetUserByUID(uidStr)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("please first login")
	}

	// Update address and sign
	err = user.UpdateAddressAndSign(form.WalletAddress, form.Sign)
	if err != nil {
		return nil, err
	}

	loginResponse := u.GenerateLoginResponse(form.WalletAddress, uidStr)

	return loginResponse, nil
}

func (u *User) GetUserAddressAndSign(UID uuid.UUID) (*LoginResponse, error) {
	uidStr := UID.String()

	// 获取用户
	user, err := u.GetUserByUID(uidStr)
	if err != nil {
		return nil, err
	}

	// 检查用户是否存在
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 创建登录响应
	loginResponse := &LoginResponse{
		WalletAddress: user.Address,
		Sign:          user.Sign,
	}

	return loginResponse, nil
}

func (u *User) UnlinkWallet(UID uuid.UUID) error {
	uidStr := UID.String()

	// 获取用户
	user, err := u.GetUserByUID(uidStr)
	if err != nil {
		return err
	}

	// 检查用户是否存在
	if user == nil {
		return errors.New("user not found")
	}

	// 更新用户对象的字段
	user.Address = ""
	user.Sign = ""

	// 保存用户
	err = user.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAppPublishAmount(uid string, amount int64) error {
	db := initializers.DB
	result := db.Model(&User{}).Where("id = ?", uid).Update("app_publish_amount", amount)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateAppAmount(uid string, amount int64) error {
	db := initializers.DB
	result := db.Model(&User{}).Where("id = ?", uid).Update("app_amount", amount)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func IsAppPublishAmountGreaterThanTen(uid string) (bool, error) {
	config, err := initializers.LoadEnv(".")
	if err != nil {
		logging.Danger("🚀 Could not load environment variables %s", err.Error())
		return false, err
	}

	db := initializers.DB
	var user User

	// 查询用户的 AppPublishAmount
	result := db.Select("app_publish_amount").Where("id = ?", uid).First(&user)
	if result.Error != nil {
		return false, result.Error
	}

	// 判断是否大于 10
	return user.AppPublishAmount > config.MaxPublishAppAmount, nil
}

func (u *User) PreCheckBalance() bool {
	if u.Balance < 1 {
		return false
	}
	return true
}

func UpdateUserFrom(id string, from string) error {
	db := initializers.DB

	// 查询用户
	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return err
	}

	// 更新用户信息
	user.FromCode = from

	// 保存更新后的用户
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func UpdateUserCodeAmount(id string, amount int) error {
	db := initializers.DB

	// 查询用户
	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return err
	}

	// 更新用户信息
	user.UsedCodeAmount = user.UsedCodeAmount + amount

	// 保存更新后的用户
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
