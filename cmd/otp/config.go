package otp

import "time"

type Config struct {
	OTPExpiry   time.Duration
	OTPLength   int
	ResendDelay time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		OTPExpiry:   5 * time.Minute,
		OTPLength:   6,
		ResendDelay: 30 * time.Second,
	}
}
