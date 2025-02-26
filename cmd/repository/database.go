package repository

import "gorm.io/gorm"

type Database interface {
	Close() error
	GetDB() *gorm.DB
}
