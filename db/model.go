package db

import "time"

type Model interface {
	GetAllTodos() ([]Todo, error)
	GetManyTodos(time.Time, time.Time) ([]Todo, error)
	GetOneTodoById(int64) (Todo, error)
	UpdateOneTodo(Todo) (Todo, error)
	DeleteOneTodo(int64) error
	CreateOneTodo(CreateTodoParams) (Todo, error)
}

// Todo ORM model structure
type Todo struct {
	Id          int64     `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Completion  float32   `json:"completion" gorm:"not null"`
	Expiry      time.Time `json:"expiry" gorm:"not null"`
	IsDone      bool      `json:"is_done"`
}
