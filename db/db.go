package db

import (
	"gorm.io/gorm"
)

type Queries struct {
	db *gorm.DB
}

// New is the function used to create an instance of the data package. It returns the type
// Models, which embeds all the types we want to be available to our application.
func New(db *gorm.DB) Model {
	if db != nil {
		db.AutoMigrate(&Todo{})
	}
	return &Queries{
		db: db,
	}
}
