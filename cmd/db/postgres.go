package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresDB struct {
	DB *gorm.DB
}

func CreatePostgresDB(dsn string) (*PostgresDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	postgresqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return postgresqlDB.Close()
}

func (p *PostgresDB) GetDB() *gorm.DB {
	if p.DB == nil {
		log.Panic("Database connection is nil")
	}
	return p.DB
}
