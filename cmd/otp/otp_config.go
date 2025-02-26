package otp

import (
	"github.com/google/uuid"
	"time"
)

type OTPStatus string
type DeliveryMethod string

const (
	DeliverySMS   DeliveryMethod = "SMS"
	DeliveryEmail DeliveryMethod = "EMAIL"
)

const (
	StatusPending  OTPStatus = "PENDING"
	StatusVerified OTPStatus = "VERIFIED"
	StatusExpired  OTPStatus = "EXPIRED"
	StatusFailed   OTPStatus = "FAILED"
)

type OTP struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Purpose      string    `gorm:"type:varchar(50);not null"`
	HashedOTP    string    `gorm:"type:text;not null"`
	Delivery     string    `gorm:"type:varchar(20);not null"`
	MobileNumber string    `gorm:"type:varchar(20)"`
	Email        string    `gorm:"type:varchar(100)"`
	Payload      string    `gorm:"type:json;not null"`
	SMSContent   string    `gorm:"type:text"`
	EmailContent string    `gorm:"type:text"`
	RetryLimit   int       `gorm:"not null"`
	RetryCount   int       `gorm:"default:0"`
	ExpiresAt    time.Time `gorm:"not null"`
	Status       string    `gorm:"type:varchar(20);not null;default:'PENDING'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OTPConfig struct {
	MinLength         int      `json:"min_length" yaml:"min_length"`
	MaxLength         int      `json:"max_length" yaml:"max_length"`
	ExpirationSeconds int      `json:"expiration_seconds" yaml:"expiration_seconds"`
	RetryLimit        int      `json:"retry_limit" yaml:"retry_limit"`
	AllowedDeliveries []string `json:"allowed_delivery_methods" yaml:"allowed_delivery_methods"`
}
