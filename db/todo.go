package db

import (
	"errors"
	"time"
)

type CreateTodoParams struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Expiry      time.Time `json:"expiry"`
}

// Returns all Todos from database
func (q *Queries) GetAllTodos() ([]Todo, error) {
	var todos []Todo
	result := q.db.Find(&todos)
	return todos, result.Error
}

// Inserts single Todo to database
func (q *Queries) CreateOneTodo(params CreateTodoParams) (Todo, error) {
	todo := Todo{
		Title:       params.Title,
		Description: params.Description,
		Expiry:      params.Expiry,
		IsDone:      false,
		Completion:  0,
	}
	result := q.db.Create(&todo)
	return todo, result.Error
}

// Returns slice of unfinished Todos from database between two "yyyy-mm-dd" dates.
func (q *Queries) GetManyTodos(startDate, endDate time.Time) ([]Todo, error) {
	var todos []Todo
	result := q.db.Where("(expiry BETWEEN ? AND ?) AND NOT is_done", startDate, endDate).Find(&todos)
	return todos, result.Error
}

// Returns single Todo with given Id
func (q *Queries) GetOneTodoById(id int64) (Todo, error) {
	todo := Todo{Id: id}
	result := q.db.First(&todo)
	if result.RowsAffected == 0 {
		return todo, errors.New("not found")
	}
	return todo, result.Error
}

// Updates existing Todo
func (q *Queries) UpdateOneTodo(todo Todo) (Todo, error) {
	result := q.db.Save(&todo)
	return todo, result.Error
}

// Deletes Todo with given Id
func (q *Queries) DeleteOneTodo(id int64) error {
	result := q.db.Delete(&Todo{}, id)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return result.Error
}
