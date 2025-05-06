package models

import (
	"github.com/google/uuid"
	"time"
)

type BillingHistory struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Operator          string    `gorm:"type:varchar(255);not null" json:"operator"`
	AccountEmail      string    `gorm:"not null" json:"account_email"`
	Amount            int64     `gorm:"bigint;not null" json:"amount"`
	TransactionType   string    `gorm:"type:varchar(255)" json:"transaction_type"`
	TransactionDetail string    `gorm:"type:varchar(255)" json:"transaction_detail"`
	TransactionTime   time.Time `gorm:"not null" json:"transaction_time"`
}
