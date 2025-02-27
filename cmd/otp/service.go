package otp

import (
	"errors"
	"fmt"
	"github.com/Zaman-R/otp-validator/cmd/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type OTPService struct {
	db     *gorm.DB
	config *config.OTPConfig
}

func NewOTPService(db *gorm.DB) (*OTPService, error) {
	var otpConfig config.OTPConfig
	return &OTPService{
		db:     db,
		config: &otpConfig,
	}, nil
}

func (s *OTPService) CreateOTP(purpose string, delivery string, mobile *string, email *string, payload map[string]interface{}, smsBody *string, emailBody *string, length int) (*OTP, string, error) {
	err := ValidateOTPRequest(s.config, purpose, delivery, mobile, email, length)
	if err != nil {
		return nil, "", err
	}

	otp, rawOTP, err := GenerateOTP(purpose, DeliveryMethod(delivery), mobile, email, payload, smsBody, emailBody, length, s.config.ExpirationSeconds, s.config.RetryLimit)
	if err != nil {
		return nil, "", err
	}

	if err := s.db.Create(otp).Error; err != nil {
		return nil, "", err
	}

	return otp, rawOTP, nil
}

func (s *OTPService) VerifyOTP(id string, userOTP string) error {
	var otp OTP
	err := s.db.Where("id = ?", id).First(&otp).Error
	if err != nil {
		return errors.New("OTP not found")
	}

	// Check expiration
	if time.Now().After(otp.ExpiresAt) {
		otp.Status = string(StatusExpired)
		s.db.Save(&otp)
		return errors.New("OTP expired")
	}

	fmt.Println("Stored OTP Hash:", otp.HashedOTP)
	fmt.Println("User Input OTP:", userOTP)

	// Compare OTP
	if err := bcrypt.CompareHashAndPassword([]byte(otp.HashedOTP), []byte(userOTP)); err != nil {
		otp.RetryCount++
		if otp.RetryCount >= otp.RetryLimit {
			otp.Status = string(StatusFailed)
		}
		s.db.Save(&otp)
		return errors.New("invalid OTP")
	}

	// Mark OTP as verified
	otp.Status = string(StatusVerified)
	s.db.Save(&otp)
	return nil
}

func (s *OTPService) ResendOTP(otpID string) (*OTP, string, error) {
	var otp *OTP
	if err := s.db.Where("id = ?", otpID).First(&otp).Error; err != nil {
		return nil, "", err
	}

	log.Printf("OTP Status: %s, ExpiresAt: %v", otp.Status, otp.ExpiresAt)

	if otp.Status != string(StatusExpired) && otp.Status != string(StatusPending) {
		return nil, "", errors.New("cannot resend OTP in current state")
	}

	if otp.RetryCount >= otp.RetryLimit {
		return nil, "", errors.New("retry limit reached")
	}
	newExpiresAt := time.Now().Add(time.Duration(s.config.ExpirationSeconds) * time.Second)

	newOTP, newRawOTP, err := GenerateOTP(
		otp.Purpose,
		DeliveryMethod(otp.Delivery),
		&otp.MobileNumber,
		&otp.Email,
		nil,
		&otp.SMSContent,
		&otp.EmailContent,
		6,
		int(time.Until(newExpiresAt).Seconds()),
		otp.RetryLimit,
	)
	if err != nil {
		return nil, "", err
	}

	otp.HashedOTP = newOTP.HashedOTP
	otp.ExpiresAt = newOTP.ExpiresAt
	otp.Status = string(StatusPending)
	otp.RetryCount++
	s.db.Save(&otp)

	return otp, newRawOTP, nil
}
