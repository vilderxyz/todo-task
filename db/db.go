package db

import (
	"gorm.io/gorm"
)

var db *gorm.DB

type Models struct {
	Todo Todo
}

// New is the function used to create an instance of the data package. It returns the type
// Models, which embeds all the types we want to be available to our application.
func New(dbPool *gorm.DB) Models {
	db = dbPool
	db.AutoMigrate(&Todo{})
	return Models{
		Todo: Todo{},
	}
}
