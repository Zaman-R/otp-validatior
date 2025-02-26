package otp

import (
	"crypto/rand"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"time"
)

func IfNotNil(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func GenerateNumericOTP(length int) (string, error) {
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp[i] = '0' + byte(num.Int64())
	}
	return string(otp), nil
}

func HashOTP(otp string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func GenerateOTP(
	purpose string,
	delivery DeliveryMethod,
	mobile *string,
	email *string,
	payload map[string]interface{},
	smsBody *string,
	emailBody *string,
	length int,
	expirationSeconds int,
	retryLimit int,

) (*OTP, string, error) {

	rawOTP, err := GenerateNumericOTP(length)
	if err != nil {
		return nil, "", err
	}

	hashedOTP, err := HashOTP(rawOTP)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(expirationSeconds) * time.Second)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}

	otp := &OTP{
		ID:           uuid.New(),
		Purpose:      purpose,
		HashedOTP:    hashedOTP,
		Delivery:     string(delivery),
		MobileNumber: IfNotNil(mobile),
		Email:        IfNotNil(email),
		Payload:      string(payloadBytes),
		SMSContent:   IfNotNil(smsBody),
		EmailContent: IfNotNil(emailBody),
		CreatedAt:    now,
		ExpiresAt:    expiresAt,
		RetryLimit:   retryLimit,
	}

	return otp, rawOTP, nil
}

func generateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
