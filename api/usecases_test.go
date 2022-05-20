package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vilderxyz/todos/db"
	"github.com/vilderxyz/todos/mock"
)

// Global mock object for testing
var todo db.Todo = db.Todo{
	Id:          123,
	Title:       "title",
	Description: "desc",
	Expiry:      time.Now().Add(time.Hour),
	IsDone:      false,
	Completion:  50,
}

type CreateTodoCase struct {
	name          string
	body          gin.H
	buildStubs    func(model *mock.MockDB)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func getCreateTodoCases(t *testing.T) []CreateTodoCase {
	return []CreateTodoCase{
		{
			name: "StatusOK",
			body: gin.H{
				"title":       todo.Title,
				"description": todo.Description,
				"expiry":      "2022-05-22",
			},
			buildStubs: func(model *mock.MockDB) {
				expiryTime, err := time.Parse("2006-01-02", "2022-05-22")
				require.NoError(t, err)
				req := db.CreateTodoParams{
					Title:       todo.Title,
					Description: todo.Description,
					Expiry:      expiryTime,
				}
				model.EXPECT().
					CreateOneTodo(gomock.Eq(req)).
					Times(1).
					Return(todo, err)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest - validation error",
			body: gin.H{
				"title":  todo.Title,
				"expiry": "2022-05-22",
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					CreateOneTodo(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - invalid date values",
			body: gin.H{
				"title":       todo.Title,
				"expiry":      "2022-13-23",
				"description": todo.Description,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					CreateOneTodo(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - expiry date from the past",
			body: gin.H{
				"title":       todo.Title,
				"expiry":      "2010-05-16",
				"description": todo.Description,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					CreateOneTodo(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError - database connection lost",
			body: gin.H{
				"title":       todo.Title,
				"expiry":      "2222-05-01",
				"description": todo.Description,
			},
			buildStubs: func(model *mock.MockDB) {
				expiryTime, err := time.Parse("2006-01-02", "2222-05-01")
				require.NoError(t, err)
				req := db.CreateTodoParams{
					Title:       todo.Title,
					Description: todo.Description,
					Expiry:      expiryTime,
				}
				model.EXPECT().
					CreateOneTodo(gomock.Eq(req)).
					Times(1).
					Return(todo, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
}

type GetTodoCase struct {
	name          string
	todoId        int64
	buildStubs    func(model *mock.MockDB)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func getGetTodoCases(t *testing.T) []GetTodoCase {
	return []GetTodoCase{
		{
			name:   "StatusOK",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "BadRequest - invalid id",
			todoId: -todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, fmt.Errorf("not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError - database connection",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
}

type UpdateTodoTextCase struct {
	name          string
	body          gin.H
	buildStubs    func(model *mock.MockDB)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func getUpdateTodoTextCases(t *testing.T) []UpdateTodoTextCase {
	updatedTitle := "t"
	updatedDesc := "d"
	updatedExpiry := "2022-05-30"

	return []UpdateTodoTextCase{
		{
			name: "StatusOK",
			body: gin.H{
				"title":       updatedTitle,
				"description": updatedDesc,
				"expiry":      updatedExpiry,
				"id":          todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				expiryTime, err := time.Parse("2006-01-02", updatedExpiry)
				require.NoError(t, err)

				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.Description = updatedDesc
				todo.Title = updatedTitle
				todo.Expiry = expiryTime

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest - invalid body",
			body: gin.H{
				"title":  updatedTitle,
				"expiry": updatedExpiry,
				"id":     todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - invalid expiry date",
			body: gin.H{
				"title":       updatedTitle,
				"expiry":      "2021-01-01",
				"description": updatedDesc,
				"id":          todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - invalid date values",
			body: gin.H{
				"title":       updatedTitle,
				"expiry":      "2021-13-01",
				"description": updatedDesc,
				"id":          todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"title":       updatedTitle,
				"description": updatedDesc,
				"expiry":      updatedExpiry,
				"id":          todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				expiryTime, err := time.Parse("2006-01-02", updatedExpiry)
				require.NoError(t, err)

				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, fmt.Errorf("not found"))

				todo.Description = updatedDesc
				todo.Title = updatedTitle
				todo.Expiry = expiryTime

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError = database connection in get",
			body: gin.H{
				"title":       updatedTitle,
				"description": updatedDesc,
				"expiry":      updatedExpiry,
				"id":          todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				expiryTime, err := time.Parse("2006-01-02", updatedExpiry)
				require.NoError(t, err)

				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, sql.ErrConnDone)

				todo.Description = updatedDesc
				todo.Title = updatedTitle
				todo.Expiry = expiryTime

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError = database connection in update",
			body: gin.H{
				"title":       updatedTitle,
				"description": updatedDesc,
				"expiry":      updatedExpiry,
				"id":          todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				expiryTime, err := time.Parse("2006-01-02", updatedExpiry)
				require.NoError(t, err)

				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.Description = updatedDesc
				todo.Title = updatedTitle
				todo.Expiry = expiryTime

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(1).
					Return(todo, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
}

type UpdateTodoCompletionCase struct {
	name          string
	body          gin.H
	buildStubs    func(model *mock.MockDB)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func getUpdateTodoCompletionCases(t *testing.T) []UpdateTodoCompletionCase {
	updatedCompletion := 51

	return []UpdateTodoCompletionCase{
		{
			name: "StatusOK",
			body: gin.H{
				"completion": updatedCompletion,
				"id":         todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.Completion = float32(updatedCompletion)

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest - invalid body",
			body: gin.H{
				"id": todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(0)

				todo.Completion = float32(updatedCompletion)

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"id":         todo.Id,
				"completion": updatedCompletion,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, fmt.Errorf("not found"))

				todo.Completion = float32(updatedCompletion)

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError - database connection in get",
			body: gin.H{
				"id":         todo.Id,
				"completion": updatedCompletion,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, sql.ErrConnDone)

				todo.Completion = float32(updatedCompletion)

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest - completion lower than value in database",
			body: gin.H{
				"id":         todo.Id,
				"completion": updatedCompletion,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.Completion = float32(updatedCompletion)

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError - database connection in update",
			body: gin.H{
				"id":         todo.Id,
				"completion": updatedCompletion,
			},
			buildStubs: func(model *mock.MockDB) {
				todo.Completion = 0

				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.Completion = float32(updatedCompletion)

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(1).
					Return(todo, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
}

type UpdateTodoDoneCase struct {
	name          string
	body          gin.H
	buildStubs    func(model *mock.MockDB)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func getUpdateTodoDoneCases(t *testing.T) []UpdateTodoDoneCase {

	return []UpdateTodoDoneCase{
		{
			name: "StatusOK",
			body: gin.H{
				"is_done": true,
				"id":      todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.IsDone = true

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest - invalid body",
			body: gin.H{
				"id": todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(0)

				todo.IsDone = true

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"is_done": true,
				"id":      todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, fmt.Errorf("not found"))

				todo.IsDone = true

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError - database connection in get",
			body: gin.H{
				"is_done": true,
				"id":      todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, sql.ErrConnDone)

				todo.IsDone = true

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest - value not affected",
			body: gin.H{
				"is_done": true,
				"id":      todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				todo.IsDone = true
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.IsDone = true

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError - database connection in update",
			body: gin.H{
				"is_done": true,
				"id":      todo.Id,
			},
			buildStubs: func(model *mock.MockDB) {
				todo.IsDone = false
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)

				todo.IsDone = true

				model.EXPECT().
					UpdateOneTodo(gomock.Eq(todo)).
					Times(1).
					Return(todo, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

}

type DeleteTodoCase struct {
	name          string
	todoId        int64
	buildStubs    func(model *mock.MockDB)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func getDeleteTodoCases(t *testing.T) []DeleteTodoCase {

	return []DeleteTodoCase{
		{
			name:   "StatusOK",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					DeleteOneTodo(gomock.Eq(todo.Id)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "BadRequest - invalid id",
			todoId: -todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					DeleteOneTodo(gomock.Eq(todo.Id)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					DeleteOneTodo(gomock.Eq(todo.Id)).
					Times(1).
					Return(fmt.Errorf("not found"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError - database connection",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					DeleteOneTodo(gomock.Eq(todo.Id)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

}

type GetTodosCase struct {
	name          string
	period        string
	buildStubs    func(model *mock.MockDB)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func getGetTodosCases(t *testing.T) []GetTodosCase {
	todos := []db.Todo{
		todo,
	}

	return []GetTodosCase{
		{
			name:   "StatusOK - get all todos",
			period: "",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetAllTodos().
					Times(1).
					Return(todos, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			period: "????",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetAllTodos().
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError - get all database connection",
			period: "",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetAllTodos().
					Times(1).
					Return(todos, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "StatusOK - get all todos for today",
			period: "today",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetManyTodos(gomock.Any(), gomock.Any()).
					Times(1).
					Return(todos, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InternalError - get all todos for today database connection",
			period: "today",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetManyTodos(gomock.Any(), gomock.Any()).
					Times(1).
					Return(todos, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "StatusOK - get all todos for tomorrow",
			period: "tomorrow",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetManyTodos(gomock.Any(), gomock.Any()).
					Times(1).
					Return(todos, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InternalError - get all todos for tomorrow database connection",
			period: "tomorrow",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetManyTodos(gomock.Any(), gomock.Any()).
					Times(1).
					Return(todos, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "StatusOK - get all todos for week",
			period: "week",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetManyTodos(gomock.Any(), gomock.Any()).
					Times(1).
					Return(todos, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InternalError - get all todos for week database connection",
			period: "week",
			buildStubs: func(model *mock.MockDB) {
				model.EXPECT().
					GetManyTodos(gomock.Any(), gomock.Any()).
					Times(1).
					Return(todos, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

}
