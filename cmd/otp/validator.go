package otp

import (
	"errors"
	"github.com/Zaman-R/otp-validator/cmd/config"
)

func ValidateOTPRequest(config *config.OTPConfig, purpose, delivery string, mobile, email *string, length int) error {
	if purpose == "" {
		return errors.New("purpose cannot be empty")
	}

	isValidDelivery := false
	for _, d := range config.AllowedDeliveries {
		if delivery == d {
			isValidDelivery = true
			break
		}
	}
	if !isValidDelivery {
		return errors.New("invalid delivery method")
	}

	if delivery == string(DeliverySMS) && (mobile == nil || *mobile == "") {
		return errors.New("mobile number is required for SMS OTP")
	}
	if delivery == string(DeliveryEmail) && (email == nil || *email == "") {
		return errors.New("email is required for email OTP")
	}

	if length < config.MinLength || length > config.MaxLength {
		return errors.New("OTP length must be within allowed range")
	}

	if config.ExpirationSeconds <= 0 {
		return errors.New("expiration time must be greater than 0")
	}
	if config.RetryLimit < 1 {
		return errors.New("retry limit must be at least 1")
	}

	return nil
}
