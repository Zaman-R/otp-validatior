package client

import (
	"errors"
	"fmt"
)

type SMSProvider interface {
	//SendSMS(phone string, message string) error
	SendSMS(phone, otp string) error
}

type CustomSMS struct{}

func (c *CustomSMS) SendSMS(phone, otp string) error {
	fmt.Printf("Sending SMS to %s: Your OTP is %s\n", phone, otp)
	return errors.New("not Implemented")
}
