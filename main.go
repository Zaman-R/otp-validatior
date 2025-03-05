package main

import (
	"fmt"
	"github.com/Zaman-R/otp-validator/cmd/client"
	"github.com/Zaman-R/otp-validator/cmd/repository"
	"log"

	"github.com/Zaman-R/otp-validator/cmd/config"
	"github.com/Zaman-R/otp-validator/cmd/otp"
)

//✅ Database initializes properly
//✅ OTP repository is correctly created
//✅ OTP service is instantiated with valid arguments

func main() {
	// Load configurations
	config.LoadConfig()
	config.ConnectDB()
	db := config.GetDB()

	// Create OTP repository
	otpRepo := repository.NewOTPRepository(db.GetDB())

	// Initialize providers (Clients can implement their own)
	smsProvider := client.NewCustomSMSProvider()
	emailProvider := client.NewCustomEmailProvider()

	// Initialize OTP Service
	otpService := otp.NewOTPService(otpRepo, smsProvider, emailProvider)

	// Example: Sending an OTP
	otpRef, err := otpService.SendOTP(otp.SendOTPRequest{
		MobileNumber: strPtr("+123456789"),
		Length:       6,
		RetryLimit:   3,
		Expiration:   300, // 5 minutes
	})

	if err != nil {
		log.Fatalf("Failed to send OTP: %v", err)
	}

	fmt.Println("OTP Sent! Reference:", otpRef)
}

// Helper function to get a string pointer
func strPtr(s string) *string {
	return &s
}
