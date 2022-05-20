package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createTodo(t *testing.T) Todo {
	todo, err := testQueries.CreateOneTodo(CreateTodoParams{
		Title:       "test_title",
		Description: "test_desc",
		Expiry:      time.Now(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, todo)
	return todo
}

func TestCreateTodo(t *testing.T) {
	createTodo(t)
}

func TestGetMany(t *testing.T) {
	for i := 0; i < 10; i++ {
		_ = createTodo(t)
	}

	todos, err := testQueries.GetAllTodos()
	require.NoError(t, err)
	require.NotEmpty(t, todos)
	require.GreaterOrEqual(t, len(todos), 10)

	for _, todo := range todos {
		require.NotEmpty(t, todo)
	}
}

func TestUpdateTodo(t *testing.T) {
	todo := createTodo(t)

	todo.Title = "New title"
	todo.Description = "New desc"
	todo.Expiry = time.Now().Add(time.Hour)
	todo.Completion = 21.37
	todo.IsDone = true

	updatedTodo, err := testQueries.UpdateOneTodo(todo)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTodo)

	require.Equal(t, todo.Id, updatedTodo.Id)
	require.Equal(t, todo.Title, updatedTodo.Title)
	require.Equal(t, todo.Description, updatedTodo.Description)
	require.Equal(t, todo.IsDone, updatedTodo.IsDone)
	require.Equal(t, todo.Completion, updatedTodo.Completion)
	require.WithinDuration(t, todo.Expiry, updatedTodo.Expiry, time.Second)
}

func TestDeleteTodo(t *testing.T) {
	todo := createTodo(t)

	err := testQueries.DeleteOneTodo(todo.Id)
	require.NoError(t, err)

	err = testQueries.DeleteOneTodo(todo.Id)
	require.Error(t, err)
}

func TestGetOneById(t *testing.T) {
	todo := createTodo(t)

	recievedTodo, err := testQueries.GetOneTodoById(todo.Id)
	require.NoError(t, err)
	require.NotEmpty(t, recievedTodo)

	require.Equal(t, todo.Id, recievedTodo.Id)
	require.Equal(t, todo.Title, recievedTodo.Title)
	require.Equal(t, todo.Description, recievedTodo.Description)
	require.Equal(t, todo.Completion, recievedTodo.Completion)
	require.Equal(t, todo.IsDone, recievedTodo.IsDone)
	require.WithinDuration(t, todo.Expiry, recievedTodo.Expiry, time.Second)

	err = testQueries.DeleteOneTodo(todo.Id)
	require.NoError(t, err)

	recievedTodo, err = testQueries.GetOneTodoById(todo.Id)
	require.Error(t, err)

}

func TestGetManyTodos(t *testing.T) {
	todo1 := createTodo(t)
	todo2 := createTodo(t)

	todo1.Expiry = time.Now().AddDate(0, 0, 1)
	todo2.Expiry = time.Now().AddDate(0, 0, 2)

	_, err := testQueries.UpdateOneTodo(todo1)
	require.NoError(t, err)

	_, err = testQueries.UpdateOneTodo(todo2)
	require.NoError(t, err)

	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, 5)
	todos, err := testQueries.GetManyTodos(startDate, endDate)
	require.NoError(t, err)
	require.NotEmpty(t, todos)
	require.Greater(t, len(todos), 0)

}
