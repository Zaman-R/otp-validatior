package config

import (
	"fmt"
	"github.com/Zaman-R/otp-validator/cmd/db"
	"github.com/Zaman-R/otp-validator/cmd/repository"
	"log"
)

var database repository.Database

func getDsn() string {
	if AppConfig == nil {
		log.Panic("❌ AppConfig is not initialized. Call config.LoadConfig() first.")
	}

	switch AppConfig.DBDriver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			AppConfig.DBHost, AppConfig.DBUser, AppConfig.DBPassword,
			AppConfig.DBName, AppConfig.DBPort, AppConfig.SSLMode,
			AppConfig.TimeZone,
		)

	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBHost,
			AppConfig.DBPort, AppConfig.DBName,
		)

	default:
		log.Panic("❌ Unsupported database driver:", AppConfig.DBDriver)
		return ""
	}
}

func ConnectDB() {
	if AppConfig == nil {
		log.Panic("❌ AppConfig is not initialized. Call config.LoadConfig() first.")
	}

	var err error
	switch AppConfig.DBDriver {
	case "postgres":
		database, err = db.CreatePostgresDB(getDsn())
	case "mysql":
		database, err = db.CreateMySQLDB(getDsn())
	default:
		log.Panic("❌ Unsupported database driver:", AppConfig.DBDriver)
	}

	if err != nil {
		log.Panic("❌ Failed to connect to database:", err)
	}

	log.Println("✅ Database connected successfully.")
}

func GetDB() repository.Database {
	if database == nil {
		log.Panic("❌ Database connection is not initialized. Call ConnectDB() first.")
	}
	return database
}
