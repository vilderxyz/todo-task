package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vilderxyz/todos/db"
)

// General response object for successful requests.
//
// Data field can be omitted.
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Request object for createTodo.
//
// Title and Description should have a minimum 1 character.
//
// Expiry must be a future date and in given format "yyyy-mm-dd".
//
// Otherwise throws 400 status.
//
// Example:
// 	{
//		"title": 		 "Clean house"
//		"description":	"I need to clean my house till 2022-12-23"
//		"expiry":		 "2022-12-23"
//	}
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,min=1"`
	Description string `json:"description" binding:"required,min=1"`
	Expiry      string `json:"expiry" binding:"required" time_format:"2006-01-02"`
}

// Validates request body and stores new Todo object in database.
func (s *Server) createTodo(ctx *gin.Context) {
	req := CreateTodoRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	expiryTime, err := time.Parse("2006-01-02", req.Expiry)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if expiryTime.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("wrong date")))
		return
	}

	res, err := s.Queries.CreateOneTodo(db.CreateTodoParams{
		Title:       req.Title,
		Description: req.Description,
		Expiry:      expiryTime,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Created todo",
		Data:    res,
	})
}

// Request object that must contain uri with Id.
//
// Id must be greater then 1.
//
// Example:
//	"http://localhost/todos/Id"
type GetTodoByIdRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

// Returns Todo object for given Id.
//
// Throws 404 status when not found.
func (s *Server) getTodoById(ctx *gin.Context) {
	req := GetTodoByIdRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res, err := s.Queries.GetOneTodoById(req.Id)
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Recieved todo",
		Data:    res,
	})
}

// Request object for updateTodoTextInfo.
//
// Title and Description should have a minimum 1 character.
//
// Expiry must be a future date and in given format "yyyy-mm-dd".
//
// Id must be greater then 1.
//
// Otherwise throws 400 status.
//
// Example:
// 	{
//		"id":			 123
//		"title": 		 "Clean house"
//		"description":	"I need to clean my house till 2022-12-23"
//		"expiry":		 "2022-12-23"
//	}
type UpdateTodoInfoRequest struct {
	Id          int64  `json:"id" binding:"required,min=1"`
	Title       string `json:"title" binding:"required,min=1"`
	Description string `json:"description" binding:"required,min=1"`
	Expiry      string `json:"expiry" binding:"required" time_format:"2006-01-02"`
}

// Finds Todo object from database for given Id. Throws 404 when not found.
//
// Then it replaces its Title, Description and Expiry parameters
// with those from request and stores updated object back to the database.
func (s *Server) updateTodoTextInfo(ctx *gin.Context) {
	req := UpdateTodoInfoRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	expiryTime, err := time.Parse("2006-01-02", req.Expiry)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if expiryTime.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("wrong date")))
		return
	}

	todo, err := s.Queries.GetOneTodoById(req.Id)
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	todo.Description = req.Description
	todo.Expiry = expiryTime
	todo.Title = req.Title

	res, err := s.Queries.UpdateOneTodo(todo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Updated todo's description",
		Data:    res,
	})
}

// Request object for updateTodoCompletionInfo.
//
// Completion needs to be between 0 and 100.
//
// Id must be greater then 1.
//
// Otherwise throws 400 status.
//
// Example:
// 	{
//		"id":			 123
//		"completion":	 99.99
//		"expiry":		 "2022-12-23"
//	}
type UpdateTodoCompletionRequest struct {
	Id         int64   `json:"id" binding:"required,min=1"`
	Completion float32 `json:"completion" binding:"required,gte=0,lte=100"`
}

// Finds Todo object from database for given Id. Throws 404 status when not found.
//
// Then it replaces its Completion parameter with requested one and stores it back in database.
//
// It throws 400 status when requested completion value is lower than the actual one.
func (s *Server) updateTodoCompletionInfo(ctx *gin.Context) {
	req := UpdateTodoCompletionRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	todo, err := s.Queries.GetOneTodoById(req.Id)
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if todo.Completion >= req.Completion {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("requsted completion progress is lower then the actual one")))
		return
	}

	todo.Completion = req.Completion

	res, err := s.Queries.UpdateOneTodo(todo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Updated todo's completion progress",
		Data:    res,
	})
}

// Request object for updateTodoCompletionInfo.
//
// IsDone must be true.
//
// Id must be greater then 1.
//
// Otherwise throws 400 status.
//
// Example:
// 	{
//		"id":		 123
//		"is_done":	true
//	}
type UpdateTodoDoneRequest struct {
	Id     int64 `json:"id" binding:"required,min=1"`
	IsDone bool  `json:"is_done" binding:"required"`
}

// Finds Todo object from database for given Id. Throws 404 status when not found.
//
// Then it replaces its IsDone parameter with requested one and stores it back in database.
//
// It throws 400 status when Todo is already finished.
func (s *Server) updateTodoDoneInfo(ctx *gin.Context) {
	req := UpdateTodoDoneRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	todo, err := s.Queries.GetOneTodoById(req.Id)
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if todo.IsDone || !req.IsDone {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("todo is already done")))
		return
	}

	todo.IsDone = req.IsDone

	res, err := s.Queries.UpdateOneTodo(todo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Updated todo's status",
		Data:    res,
	})
}

// Request object that must contain uri with Id.
//
// Id must be greater then 1.
//
// Example:
//	"http://localhost/todos/Id"
type DeleteTodoRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

// Deletes Todo with given Id.
//
// Throws 404 when it deleted nothing.
func (s *Server) deleteTodo(ctx *gin.Context) {
	req := DeleteTodoRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.Queries.DeleteOneTodo(req.Id)
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Deleted todo",
	})
}

// Request object with period query. Can be omitted.
//
// Period must be string and of one [ "today" , "tomorrow" , "week" , ""].
//
// Otherwise throws 400 status.
//
// Examples:
//	"http://localhost/todos" 				- gets all finished and unfinished Todos
//	"http://localhost/todos?period=today"    - gets all unfinished Todos that expires after today
//	"http://localhost/todos?period=tomorrow" - gets all unfinished Todos that expires after tomorrow
//	"http://localhost/todos?period=week" 	- gets all unfinished Todos that expires after Sunday this week
type GetTodosRequest struct {
	Period string `form:"period" binding:"period"`
}

// Gets slice of Todo objects depending on given Period query.
func (s *Server) getTodos(ctx *gin.Context) {
	req := GetTodosRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var todos []db.Todo
	var err error
	var message string

	switch req.Period {
	case "today":
		day := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		endTime, err := time.Parse("2006-01-02", day)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		todos, err = s.Queries.GetManyTodos(time.Now(), endTime)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		message = "Got all todos for today"

	case "tomorrow":
		day := time.Now().AddDate(0, 0, 2).Format("2006-01-02")
		endTime, err := time.Parse("2006-01-02", day)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		todos, err = s.Queries.GetManyTodos(time.Now(), endTime)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		message = "Got all todos for tomorrow"

	case "week":
		addDays := 8 - int(time.Now().Weekday())
		day := time.Now().AddDate(0, 0, addDays).Format("2006-01-02")
		endTime, err := time.Parse("2006-01-02", day)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		todos, err = s.Queries.GetManyTodos(time.Now(), endTime)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		message = "Got all todos for this week"

	case "":
		todos, err = s.Queries.GetAllTodos()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		message = "Got all todos"

	default:
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("wrong param")))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: message,
		Data:    todos,
	})
}
