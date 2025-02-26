package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type MySQLDB struct {
	DB *gorm.DB
}

func CreateMySQLDB(dsn string) (*MySQLDB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &MySQLDB{DB: db}, nil
}

func (m *MySQLDB) Close() error {
	mysqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return mysqlDB.Close()
}

func (m *MySQLDB) GetDB() *gorm.DB {
	if m.DB == nil {
		log.Panic("Database connection is nil")
	}
	return m.DB
}
