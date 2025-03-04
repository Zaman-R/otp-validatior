package client

type EmailProvider interface {
	SendEmail(to string, subject string, body string) error
}

type CustomEmail struct{}

func (c *CustomEmail) SendEmail(to string, subject string, body string) error {
	return nil
}
