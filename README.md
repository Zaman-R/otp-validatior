# OTP Validator (Go Package) 📲

## Overview
`otp-validator` is a Go package that provides a robust **OTP (One-Time Password) authentication system** with SMS and Email support. It allows easy **OTP generation, validation, and expiration handling**, making it ideal for **user authentication**, **transaction verification**, and **multi-factor authentication (MFA)**.

## Features
- ✅ OTP Generation & Validation
- ✅ Secure Database Storage for OTPs
- ✅ Configurable OTP Expiry Time
- ✅ Supports SMS and Email-based OTP Delivery
- ✅ Customizable Storage and Notification Providers
- ✅ Rate Limiting & Expiry Management

---

## Installation

1. Install the package using `go get`:
   ```sh
   go get github.com/Zaman-R/otp-validator
   ```

2. Import the package into your project:
   ```go
   import "github.com/Zaman-R/otp-validator/cmd/otp"

   ```

---

## Configuration

### 1. Environment Variables Setup
Set up the following **environment variables** in a `.env` file or system environment:

```sh
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=yourdbname
SSL_MODE=disable
TIMEZONE=UTC
OTP_EXPIRY=5m
TOTP_ENABLED=false
```

- `OTP_EXPIRY`: Sets OTP expiration time (e.g., 5m for 5 minutes).
- `TOTP_ENABLED`: Enables **Time-based OTPs** (default: `false`).

### 2. Database Setup
Ensure your **PostgreSQL/MySQL database** is set up before running migrations.

Run database migrations:
```sh
go run cmd/db/migrate.go
```

---

## Usage

### 1. Initializing OTP Service
In `main.go`, **initialize the OTP service**:

```go
package main

import (
	"log"

	"github.com/Zaman-R/otp-validator/cmd/config"
	"github.com/Zaman-R/otp-validator/cmd/otp"
	"github.com/Zaman-R/otp-validator/cmd/repository"
	"github.com/Zaman-R/otp-validator/cmd/client"
)

func main() {
	// Load Config and Connect to Database
	config.LoadConfig()
	config.ConnectDB()
	db := config.GetDB()

	// Create OTP repository
	otpRepo := repository.NewOTPRepository(db.GetDB())

	// Initialize Custom SMS and Email Providers
	smsProvider := client.NewCustomSMSProvider()
	emailProvider := client.NewCustomEmailProvider()

	// Initialize OTP Service
	otpService := otp.NewOTPService(otpRepo, smsProvider, emailProvider)

	// Example Usage
	otpExample(otpService)
}

func otpExample(otpService *otp.OTPService) {
	// Generate OTP for User
	otpRequest := otp.OTPRequest{
		MobileNumber: "1234567890",
		Email:        "user@example.com",
	}
	otpCode, err := otpService.GenerateOTP(otpRequest)
	if err != nil {
		log.Fatal("Failed to generate OTP:", err)
	}
	log.Println("Generated OTP:", otpCode)

	// Validate OTP
	isValid := otpService.ValidateOTP("1234567890", otpCode)
	if isValid {
		log.Println("✅ OTP is valid!")
	} else {
		log.Println("❌ OTP is invalid or expired.")
	}
}
```

---

## Implementation Guide

### 2. Generating an OTP
To generate an OTP and send it via **SMS or Email**, use:

```go
otpRequest := otp.OTPRequest{
	MobileNumber: "1234567890",
	Email:        "user@example.com",
}
otpCode, err := otpService.GenerateOTP(otpRequest)
if err != nil {
	log.Fatal("❌ OTP generation failed:", err)
}
log.Println("✅ OTP sent successfully:", otpCode)
```

### 3. Validating an OTP
To validate an OTP entered by the user:

```go
isValid := otpService.ValidateOTP("1234567890", otpCode)
if isValid {
	log.Println("✅ OTP is correct!")
} else {
	log.Println("❌ OTP is incorrect or expired.")
}
```

### 4. Expiring an OTP Before Timeout
If you want to **manually expire an OTP**:

```go
otpService.ExpireOTP("1234567890")
log.Println("✅ OTP expired manually.")
```

---

## Customizing Providers (SMS & Email)
You can implement **custom SMS and Email providers** to send OTPs.

### 1. Implementing a Custom SMS Provider
Create your own **SMS sending logic**:

```go
package client

import "log"

type CustomSMSProvider struct{}

func NewCustomSMSProvider() *CustomSMSProvider {
	return &CustomSMSProvider{}
}

func (s *CustomSMSProvider) SendSMS(to string, message string) error {
	log.Printf("📩 Sending SMS to %s: %s\n", to, message)
	return nil // Replace with real SMS API call
}
```

### 2. Implementing a Custom Email Provider
Customize **Email notifications**:

```go
package client

import "log"

type CustomEmailProvider struct{}

func NewCustomEmailProvider() *CustomEmailProvider {
	return &CustomEmailProvider{}
}

func (e *CustomEmailProvider) SendEmail(to string, subject string, body string) error {
	log.Printf("📧 Sending Email to %s: %s\n", to, subject)
	return nil // Replace with actual email API
}
```

Then, **use them** in `main.go`:
```go
smsProvider := client.NewCustomSMSProvider()
emailProvider := client.NewCustomEmailProvider()
otpService := otp.NewOTPService(otpRepo, smsProvider, emailProvider)
```

---

## Database Schema
This package uses a **relational database (PostgreSQL/MySQL)** for OTP storage.

### 1. OTP Table Schema
```sql
CREATE TABLE otp_requests (
    id SERIAL PRIMARY KEY,
    mobile_number VARCHAR(20),
    email VARCHAR(255),
    otp_code VARCHAR(6),
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_used BOOLEAN DEFAULT FALSE
);
```

---

## Security Best Practices
- **Use secure OTP lengths** (6+ digits).
- **Expire OTPs quickly** (1-5 minutes recommended).
- **Rate limit OTP requests** to prevent abuse.
- **Hash OTPs** before storing them in the database.
- **Use HTTPS** for secure transmission.

---

## Troubleshooting

### ❌ Database Connection Fails
**Solution**: Ensure your database is running and credentials in `.env` are correct.

### ❌ OTP Not Sending
**Solution**: Verify SMS/Email provider implementation. Try logging messages before sending.

### ❌ OTP Always Invalid
**Solution**: Check if OTPs are stored in the database and have not expired.

---

## Contributing
We welcome contributions! 🚀  
Feel free to submit PRs or issues.

---

## License
MIT License © 2025 Zaman-R

