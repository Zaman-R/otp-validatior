package di

import (
	"log"

	"github.com/Zaman-R/otp-validator/cmd/config"
	"github.com/Zaman-R/otp-validator/cmd/otp"
	"github.com/Zaman-R/otp-validator/cmd/totp"
)

var GlobalRegistry *Registry

type Registry struct {
	OTPService  *otp.OTPService
	TOTPService *totp.TOTPService
}

func init() {
	GlobalRegistry = NewRegistry()
}

func NewRegistry() *Registry {
	db := config.GetDB()
	var otpService *otp.OTPService
	var totpService *totp.TOTPService

	if !config.TOTPEnabled() {
		var err error
		otpService, err = otp.NewOTPService(db.GetDB())
		if err != nil {
			log.Panic("❌ Failed to initialize OTP service:", err)
		}
		log.Println("✅ OTP service initialized")
	}

	return &Registry{
		OTPService:  otpService,
		TOTPService: totpService,
	}
}
