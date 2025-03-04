package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPService struct {
	Issuer      string
	AccountName string
	Secret      string
}

func NewTOTPService(issuer, accountName, existingSecret string) (*TOTPService, error) {
	var secret string
	var err error

	// If secret is already provided (e.g., during login), use it
	if existingSecret != "" {
		secret = existingSecret
	} else {
		secret, err = generateSecret()
		if err != nil {
			return nil, err
		}
	}

	return &TOTPService{
		Issuer:      issuer,
		AccountName: accountName,
		Secret:      secret,
	}, nil
}

func generateSecret() (string, error) {
	key := make([]byte, 10)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(key), nil
}

func (s *TOTPService) GenerateTOTP() (string, error) {
	code, err := totp.GenerateCodeCustom(s.Secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *TOTPService) VerifyTOTP(code string) bool {
	return totp.Validate(code, s.Secret)
}

func (s *TOTPService) GenerateTOTPURL() (string, error) {
	key, err := otp.NewKeyFromURL(fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		s.Issuer, s.AccountName, s.Secret, s.Issuer))
	if err != nil {
		return "", err
	}
	return key.URL(), nil
}

func GenerateHOTP(secret string, counter uint64) (string, error) {
	hmacHash := hmac.New(sha1.New, []byte(secret))
	_, err := hmacHash.Write([]byte(fmt.Sprintf("%d", counter)))
	if err != nil {
		return "", err
	}
	hash := hmacHash.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	binary := (int(hash[offset])&0x7F)<<24 |
		(int(hash[offset+1])&0xFF)<<16 |
		(int(hash[offset+2])&0xFF)<<8 |
		(int(hash[offset+3]) & 0xFF)

	otp := binary % int(math.Pow10(6))
	return fmt.Sprintf("%06d", otp), nil
}
