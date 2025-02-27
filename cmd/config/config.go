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
}

type TOTPConfig struct {
	Issuer     string `json:"issuer" yaml:"issuer"`
	Digits     int    `json:"digits" yaml:"digits"`
	Period     int    `json:"period" yaml:"period"`
	Skew       int    `json:"skew" yaml:"skew"`
	SecretSize int    `json:"secret_size" yaml:"secret_size"`
	Algorithm  string `json:"algorithm" yaml:"algorithm"`
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
var ConfigTOTP *TOTPConfig
var isTOTPEnabled bool

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
	}

	ConfigTOTP = &TOTPConfig{
		Issuer:     viper.GetString("TOTP_ISSUER"),
		Digits:     viper.GetInt("TOTP_DIGITS"),
		Period:     viper.GetInt("TOTP_PERIOD"),
		Skew:       viper.GetInt("TOTP_SKEW"),
		SecretSize: viper.GetInt("TOTP_SECRET_SIZE"),
		Algorithm:  viper.GetString("TOTP_ALGORITHM"),
	}

	isTOTPEnabled = viper.GetBool("TOTP_ENABLED")

	log.Println("✅ Configuration loaded successfully")
}

func TOTPEnabled() bool {
	return isTOTPEnabled
}
