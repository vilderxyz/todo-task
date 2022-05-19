package db

import (
	"errors"
	"time"
)

// Todo ORM model structure
type Todo struct {
	Id          int64     `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Completion  float32   `json:"completion" gorm:"not null"`
	Expiry      time.Time `json:"expiry" gorm:"not null"`
	IsDone      bool      `json:"is_done"`
}

// Returns all Todos from database
func (t *Todo) GetAll() ([]Todo, error) {
	var todos []Todo
	result := db.Find(&todos)
	return todos, result.Error
}

// Inserts single Todo to database
func (t *Todo) CreateOne(todo Todo) (*Todo, error) {
	result := db.Create(&todo)
	return &todo, result.Error
}

// Returns slice of unfinished Todos from database between two "yyyy-mm-dd" dates.
func (t *Todo) GetMany(startDate, endDate string) ([]Todo, error) {
	var todos []Todo
	result := db.Where("(expiry BETWEEN ? AND ?) AND NOT is_done", startDate, endDate).Find(&todos)
	return todos, result.Error
}

// Returns single Todo with given Id
func (t *Todo) GetOneById(id int64) (*Todo, error) {
	todo := Todo{Id: id}
	result := db.First(&todo)
	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return &todo, result.Error
}

// Updates existing Todo
func (t *Todo) UpdateOne(todo Todo) (*Todo, error) {
	result := db.Save(&todo)
	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return &todo, result.Error
}

// Deletes Todo with given Id
func (t *Todo) DeleteOne(id int64) error {
	result := db.Delete(&Todo{}, id)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return result.Error
}
