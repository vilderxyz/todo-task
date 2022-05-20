package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vilderxyz/todos/mock"
)

func TestCreateTodo(t *testing.T) {

	testCases := getCreateTodoCases(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockDB(ctrl)
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

	testCases := getGetTodoCases(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockDB(ctrl)
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

	testCases := getUpdateTodoTextCases(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockDB(ctrl)
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

	testCases := getUpdateTodoCompletionCases(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockDB(ctrl)
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

	testCases := getUpdateTodoDoneCases(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockDB(ctrl)
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

	testCases := getDeleteTodoCases(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockDB(ctrl)
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

	testCases := getGetTodosCases(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := mock.NewMockDB(ctrl)
			tc.buildStubs(model)

			server := newTestServer(t, model)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/todos?period=%v", tc.period)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
