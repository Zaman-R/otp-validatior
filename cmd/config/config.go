package config

import (
	"log"

	"github.com/spf13/viper"
)

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

func LoadEnv() {
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

	log.Println("✅ Configuration loaded successfully")
}
