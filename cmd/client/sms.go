package client

import (
	"errors"
	"fmt"
)

type SMSProvider interface {
	//SendSMS(phone string, message string) error
	SendSMS(phone, otp string) error
}

type CustomSMSProvider struct{}

func NewCustomSMSProvider() *CustomSMSProvider {
	return &CustomSMSProvider{}
}
func (c *CustomSMSProvider) SendSMS(phone, otp string) error {
	fmt.Printf("Sending SMS to %s: Your OTP is %s\n", phone, otp)
	return errors.New("not Implemented")
}
