package totp

import (
	"log"

	"github.com/pquerna/otp/totp"
)

func ValidateTOTP(secret, code string) bool {
	valid := totp.Validate(code, secret)
	if !valid {
		log.Println("❌ Invalid TOTP code")
		return false
	}
	log.Println("✅ TOTP code validated successfully")
	return true
}
