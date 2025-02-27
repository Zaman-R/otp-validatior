package totp

import (
	"github.com/Zaman-R/otp-validator/cmd/config"
	"log"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret(username string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.ConfigTOTP.Issuer,
		AccountName: username,
		SecretSize:  uint(config.ConfigTOTP.SecretSize),
	})
	if err != nil {
		log.Println("❌ Error generating TOTP secret:", err)
		return nil, err
	}
	return key, nil
}
