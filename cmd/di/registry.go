package di

import (
	"github.com/Zaman-R/otp-validator/cmd/config"
	"github.com/Zaman-R/otp-validator/cmd/otp"
	"log"
)

var GlobalRegistry *Registry

type Registry struct {
	OTPService *otp.OTPService
}

func init() {
	GlobalRegistry = NewRegistry()
}

func NewRegistry() *Registry {
	db := config.GetDB()
	otpService, err := otp.NewOTPService(db.GetDB())
	if err != nil {
		log.Panic(err)
	}

	return &Registry{
		OTPService: otpService,
	}
}
