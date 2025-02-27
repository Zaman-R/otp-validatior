package config

import (
	"github.com/spf13/viper"
	"log"
)

type OTPConfig struct {
	MinLength         int      `json:"min_length" yaml:"min_length"`
	MaxLength         int      `json:"max_length" yaml:"max_length"`
	ExpirationSeconds int      `json:"expiration_seconds" yaml:"expiration_seconds"`
	RetryLimit        int      `json:"retry_limit" yaml:"retry_limit"`
	AllowedDeliveries []string `json:"allowed_delivery_methods" yaml:"allowed_delivery_methods"`
	EnableTOTP        bool     `json:"enable_totp" yaml:"enable_totp"`               // ✅  Enable TOTP
	TOTPSecretLength  int      `json:"totp_secret_length" yaml:"totp_secret_length"` // ✅  Secret length
	TOTPIssuer        string   `json:"totp_issuer" yaml:"totp_issuer"`               // ✅  Issuer name for QR code
	TOTPPeriod        int      `json:"totp_period" yaml:"totp_period"`               // ✅  TOTP validity period
}

type Config struct {
	DBDriver   string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	SSLMode    string
	TimeZone   string
}

var AppConfig *Config
var ConfigOTP *OTPConfig

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()
	AppConfig = &Config{
		DBDriver:   viper.GetString("DB_DRIVER"),
		DBHost:     viper.GetString("DB_HOST"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBPort:     viper.GetString("DB_PORT"),
		SSLMode:    viper.GetString("SSL_MODE"),
		TimeZone:   viper.GetString("TIME_ZONE"),
	}

	ConfigOTP = &OTPConfig{
		MinLength:         viper.GetInt("OTP_MIN_LENGTH"),
		MaxLength:         viper.GetInt("OTP_MAX_LENGTH"),
		ExpirationSeconds: viper.GetInt("OTP_EXPIRATION_SECONDS"),
		RetryLimit:        viper.GetInt("OTP_RETRY_LIMIT"),
		AllowedDeliveries: viper.GetStringSlice("OTP_ALLOWED_DELIVERY"),
		EnableTOTP:        viper.GetBool("ENABLE_TOTP"),
		TOTPSecretLength:  viper.GetInt("TOTP_SECRET_LENGTH"),
		TOTPIssuer:        viper.GetString("TOTP_ISSUER"),
		TOTPPeriod:        viper.GetInt("TOTP_PERIOD"),
	}

	log.Println("✅ Configuration loaded successfully")
}
