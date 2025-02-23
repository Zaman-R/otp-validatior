package otp

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"time"
)

type OTPConfig struct {
	Length       int
	Expiry       time.Duration
	AllowedChars string
}

var DefaultOTPConfig = OTPConfig{
	Length:       6,
	Expiry:       5 * time.Hour,
	AllowedChars: "0123456789",
}

func GenerateOTP(config OTPConfig) (string, error) {
	if config.Length <= 0 {
		return "", errors.New("invalid OTP length")
	}
	otp := make([]byte, config.Length)
	charLen := big.NewInt(int64(len(config.AllowedChars)))

	for i := range otp {
		randIndex, err := rand.Int(rand.Reader, charLen)
		if err != nil {
			return "", err
		}
		otp[i] = config.AllowedChars[randIndex.Int64()]
	}

	return string(otp), nil
}

func ParseOTPConfig(otp string) (OTPConfig, error) {
	var config OTPConfig
	if err := json.Unmarshal([]byte(otp), &config); err != nil {
		return DefaultOTPConfig, err
	}

	if config.Length == 0 {
		config.Length = DefaultOTPConfig.Length
	}
	if config.Expiry == 0 {
		config.Expiry = DefaultOTPConfig.Expiry
	}
	if config.AllowedChars == "" {
		config.AllowedChars = DefaultOTPConfig.AllowedChars
	}

	return config, nil
}
