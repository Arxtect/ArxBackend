package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"  json:"id"`
	Name           string     `gorm:"type:varchar(255);not null;unique" json:"name"`
	Email          string     `gorm:"uniqueIndex;not null"  json:"email"`
	Password       string     `gorm:"not null" json:"password"`
	Role           string     `gorm:"type:varchar(255);not null" json:"role"`
	SharedProjects []*Project `gorm:"many2many:project_members;" json:"sharedProjects"`

	// Legacy User
	Provider           string    `gorm:"not null" json:"provider"`
	Photo              string    `gorm:"type:varchar(255)" json:"photo"`
	VerificationCode   string    `json:"verification_code"`
	PasswordResetToken string    `json:"password_reset_token"`
	PasswordResetAt    time.Time `json:"password_reset_at"`
	Verified           bool      `gorm:"not null" json:"verified"`
	Balance            int64     `gorm:"not null:default:0" json:"balance"`
	CreatedAt          time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt          time.Time `gorm:"not null" json:"updated_at"`
}

type SafeUser struct {
	Base
	Name string `json:"name"`
}

func (s *SafeUser) TableName() string {
	return "users"
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	Photo           string `json:"photo,omitempty"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
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

type SubscribeRequest struct {
	RoomId string `json:"room_id"`
}
