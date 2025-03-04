package otp

import (
	"github.com/google/uuid"
	"time"
)

type OTP struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Purpose            string    `gorm:"type:varchar(50);not null"`
	HashedOTP          string    `gorm:"type:text;not null"`
	Delivery           string    `gorm:"type:varchar(20);not null"`
	MobileNumber       string    `gorm:"type:varchar(20)"`
	Email              string    `gorm:"type:varchar(100)"`
	TransactionPayload string    `gorm:"type:text"`
	RetryLimit         int       `gorm:"not null"`
	RetryCount         int       `gorm:"default:0"`
	ExpiresAt          time.Time `gorm:"not null"`
	Status             string    `gorm:"type:varchar(20);not null;default:'PENDING'"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

const (
	OTPStatusPending  = "PENDING"
	OTPStatusUsed     = "USED"
	OTPStatusExpired  = "EXPIRED"
	OTPStatusVerified = "verified"
)
