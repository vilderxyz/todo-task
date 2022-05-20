package db

import (
	"gorm.io/gorm"
)

// Implements DB interface
type Queries struct {
	db *gorm.DB
}

// Returns object that implements DB interface
//
// Meanwhile migrates all ORM models
func New(db *gorm.DB) DB {
	if db != nil {
		db.AutoMigrate(&Todo{})
	}
	return &Queries{
		db: db,
	}
}
