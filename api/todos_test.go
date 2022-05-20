package api

import (
	"bytes"
	"encoding/json"
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

func TestCreateTodo(t *testing.T) {
	todo := db.Todo{
		Id:          123,
		Title:       "title",
		Description: "desc",
		Expiry:      time.Now().Add(time.Hour),
		IsDone:      false,
		Completion:  50,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(model *mock.MockModel)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: gin.H{
				"title":       todo.Title,
				"description": todo.Description,
				"expiry":      "2022-05-22",
			},
			buildStubs: func(model *mock.MockModel) {
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
					Return(todo, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockModel(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/todos"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetTodoById(t *testing.T) {
	todo := db.Todo{
		Id:          123,
		Title:       "title",
		Description: "desc",
		Expiry:      time.Now().Add(time.Hour),
		IsDone:      false,
		Completion:  50,
	}

	testCases := []struct {
		name          string
		todoId        int64
		buildStubs    func(model *mock.MockModel)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "StatusOK",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockModel) {
				model.EXPECT().
					GetOneTodoById(gomock.Eq(todo.Id)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockModel(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/todos/%d", tc.todoId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateTextTodo(t *testing.T) {
	todo := db.Todo{
		Id:          123,
		Title:       "title",
		Description: "desc",
		Expiry:      time.Now().Add(time.Hour),
		IsDone:      false,
		Completion:  50,
	}

	updatedTitle := "t"
	updatedDesc := "d"
	updatedExpiry := "2022-05-30"

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(model *mock.MockModel)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: gin.H{
				"title":       updatedTitle,
				"description": updatedDesc,
				"expiry":      updatedExpiry,
				"id":          todo.Id,
			},
			buildStubs: func(model *mock.MockModel) {
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
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockModel(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/todos"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCompletionTodo(t *testing.T) {
	todo := db.Todo{
		Id:          123,
		Title:       "title",
		Description: "desc",
		Expiry:      time.Now().Add(time.Hour),
		IsDone:      false,
		Completion:  50,
	}

	updatedCompletion := 51

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(model *mock.MockModel)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: gin.H{
				"completion": updatedCompletion,
				"id":         todo.Id,
			},
			buildStubs: func(model *mock.MockModel) {
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
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockModel(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/todos/completion"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateDoneTodo(t *testing.T) {
	todo := db.Todo{
		Id:          123,
		Title:       "title",
		Description: "desc",
		Expiry:      time.Now().Add(time.Hour),
		IsDone:      false,
		Completion:  50,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(model *mock.MockModel)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: gin.H{
				"is_done": true,
				"id":      todo.Id,
			},
			buildStubs: func(model *mock.MockModel) {
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
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockModel(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/todos/done"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	todo := db.Todo{
		Id:          123,
		Title:       "title",
		Description: "desc",
		Expiry:      time.Now().Add(time.Hour),
		IsDone:      false,
		Completion:  50,
	}

	testCases := []struct {
		name          string
		todoId        int64
		buildStubs    func(model *mock.MockModel)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "StatusOK",
			todoId: todo.Id,
			buildStubs: func(model *mock.MockModel) {
				model.EXPECT().
					DeleteOneTodo(gomock.Eq(todo.Id)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockModel(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/todos/%d", tc.todoId)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetTodos(t *testing.T) {
	todo := db.Todo{
		Id:          123,
		Title:       "title",
		Description: "desc",
		Expiry:      time.Now().Add(time.Hour),
		IsDone:      false,
		Completion:  50,
	}

	todos := []db.Todo{
		todo,
	}

	testCases := []struct {
		name          string
		period        string
		buildStubs    func(model *mock.MockModel)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "StatusOK",
			period: "",
			buildStubs: func(model *mock.MockModel) {
				model.EXPECT().
					GetAllTodos().
					Times(1).
					Return(todos, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockModel(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/todos?period=%s", tc.period)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
