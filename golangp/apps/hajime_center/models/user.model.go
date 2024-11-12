package models

import (
	"hajime/golangp/apps/hajime_center/initializers"
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
	Verified           bool      `gorm:"not null"`
	Balance            int64     `gorm:"not null:default:0"`
	Address            string    `gorm:"type:varchar(255);default:''"`
	Status             int32     `gorm:"not null;default:1"` // Corrected type
	Code               string    `gorm:"type:varchar(255);default:''"`
	Twitter            string    `gorm:"type:varchar(255);default:''"` // Twitter
	Telegram           string    `gorm:"type:varchar(255);default:''"` // Telegram
	Discord            string    `gorm:"type:varchar(255);default:''"` // Discord
	CreatedAt          time.Time `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	Photo           string `json:"photo,omitempty"`
	Role            string `json:"role"`
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
	Balance   int64     `json:"balance,default:0"`
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

func (u *User) GetUserByAddress(address string) (*User, error) {
	db := initializers.DB
	var user User
	result := db.Where("address = ?", address).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *User) SaveUser(user *User) error {
	db := initializers.DB
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
