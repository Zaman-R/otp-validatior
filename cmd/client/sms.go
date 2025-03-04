package client

type SMSProvider interface {
	SendSMS(phone string, message string) error
}

type CustomSMS struct{}

func (c *CustomSMS) SendSMS(phone string, message string) error {
	return nil
}
