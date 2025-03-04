package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

var OTPJwtSecretKey = []byte("your-secret-key")

type OTPPayload struct {
	OTPref string    `json:"otp_ref"`
	Iat    time.Time `json:"iat"`
	Exp    time.Time `json:"exp"`
}

type Claims struct {
	OTPRef string `json:"otp_ref"`
	jwt.RegisteredClaims
}

func GenerateToken(payload map[string]interface{}, expirationSeconds int) (string, error) {
	iat := time.Now().UTC()

	// Set standard claims
	claims := jwt.MapClaims{
		"iat": iat.Unix(),
		"exp": iat.Add(time.Duration(expirationSeconds) * time.Second).Unix(),
	}

	// Add custom payload fields
	for key, value := range payload {
		claims[key] = value
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(OTPJwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return OTPJwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func HashOTPWithTime(otp string, createdAt, expiresAt time.Time) string {
	data := otp + createdAt.String() + expiresAt.String()
	hash := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func ValidateOTP(otp, providedHash string, createdAt, expiresAt time.Time) bool {
	return HashOTPWithTime(otp, createdAt, expiresAt) == providedHash
}

func EncodeBase64(payload map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func DecodeBase64(encodedStr string) (map[string]interface{}, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return nil, err
	}

	var decodedPayload map[string]interface{}
	err = json.Unmarshal(decodedBytes, &decodedPayload)
	if err != nil {
		return nil, err
	}
	return decodedPayload, nil
}

func SendSMS(phone, otp string) error {
	fmt.Printf("Sending SMS to %s: Your OTP is %s\n", phone, otp)
	return errors.New("not Implemented")
}

func SendEmail(email, otp string) error {
	fmt.Printf("Sending Email to %s: Your OTP is %s\n", email, otp)
	return errors.New("not Implemented")
}

func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func ReplacePlaceholder(template *string, otpCode string) string {
	if template == nil {
		return ""
	}
	return strings.ReplaceAll(*template, "<otp>", otpCode)
}

func GenerateSecureOTP() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b)[:6], nil
}

func HashOTP(otp string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func DetermineDeliveryMethod(email, phone string) string {
	if phone != "" {
		return "SMS"
	}
	if email != "" {
		return "EMAIL"
	}
	return "UNKNOWN"
}

func GetStringValue(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func SanitizePayload(payload map[string]interface{}) (map[string]interface{}, error) {
	if payload == nil {
		return nil, errors.New("payload is empty")
	}

	sanitizedData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var sanitizedPayload map[string]interface{}
	err = json.Unmarshal(sanitizedData, &sanitizedPayload)
	if err != nil {
		return nil, err
	}

	return sanitizedPayload, nil
}

func GetString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func GetStringPtr(m map[string]interface{}, key string) *string {
	if val, ok := m[key].(string); ok {
		return &val
	}
	return nil
}

func GetInt(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key].(int); ok {
		return val
	}
	return defaultValue
}

func GetDuration(m map[string]interface{}, key string, defaultValue time.Duration) time.Duration {
	if val, ok := m[key].(time.Duration); ok {
		return val
	}
	return defaultValue
}

func GetMap(m map[string]interface{}, key string) map[string]interface{} {
	if val, ok := m[key].(map[string]interface{}); ok {
		return val
	}
	return nil
}
