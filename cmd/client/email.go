package client

import (
	"errors"
	"fmt"
)

type EmailProvider interface {
	//SendEmail(to string, subject string, body string) error
	SendEmail(email, otp string) error
}

type CustomEmail struct{}

func (c *CustomEmail) SendEmail(email, otp string) error {
	fmt.Printf("Sending Email to %s: Your OTP is %s\n", email, otp)
	return errors.New("not Implemented")
}
