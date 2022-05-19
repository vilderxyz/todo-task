package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vilderxyz/todos/db"
)

// General response object for successful requests.
// Data field can be ignored.
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Required object for createTodo request. Another json fields will be ignored.
// Title and Description should have a minimum 1 character.
// Expiry date must be in given format "yyyy-mm-dd" or it throws an error.
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,min=1"`
	Description string `json:"description" binding:"required,min=1"`
	Expiry      string `json:"expiry" binding:"required" time_format="2006-01-02"`
}

// Validates request's body, creates Todo object and stores it in the database.
// Expiry date must be in the future or it throws ab error.
func (s *Server) createTodo(ctx *gin.Context) {
	req := CreateTodoRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	expiryTime, err := time.Parse("2006-01-02", req.Expiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if expiryTime.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("wrong date")))
		return
	}

	res, err := s.Models.Todo.CreateOne(db.Todo{
		Title:       req.Title,
		Description: req.Description,
		Expiry:      expiryTime,
		Completion:  0.0,
		IsDone:      false,
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Created todo",
		Data:    res,
	})
}

// Request's uri to be validated. Id must be greater then 1.
type GetTodoByIdRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

// Returns Todo object for given Id. Throws 404 when not found with error message.
func (s *Server) getTodoById(ctx *gin.Context) {
	req := GetTodoByIdRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res, err := s.Models.Todo.GetOneById(req.Id)
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

// Required object for updateTodoTextInfo request. Another json fields will be ignored.
// Title and Description should have a minimum 1 character. Id must be greater then 1.
// Expiry date needs to be in given format "yyyy-mm-dd" or it throws an error.
type UpdateTodoInfoRequest struct {
	Id          int64  `json:"id" binding:"required,min=1"`
	Title       string `json:"title" binding:"required,min=1"`
	Description string `json:"description" binding:"required,min=1"`
	Expiry      string `json:"expiry" binding:"required" time_format="2006-01-02"`
}

// Finds Todo object from database for given Id. Replaces its text parameters
// with those from request and stores updated object back to database.
// Expiry date must be in the future or it throws en error.
func (s *Server) updateTodoTextInfo(ctx *gin.Context) {
	req := UpdateTodoInfoRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	expiryTime, err := time.Parse("2006-01-02", req.Expiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if expiryTime.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("wrong date")))
		return
	}

	todo, err := s.Models.Todo.GetOneById(req.Id)
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res, err := s.Models.Todo.UpdateOne(db.Todo{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Expiry:      expiryTime,
		IsDone:      todo.IsDone,
		Completion:  todo.Completion,
	})
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Updated todo's description",
		Data:    res,
	})
}

// Required object for updateTodoCompletionInfo request. Another json fields are ignored.
// Id must be greater then 1. Completion needs to be between 0 and 100.
type UpdateTodoCompletionRequest struct {
	Id         int64   `json:"id" binding:"required,min=1"`
	Completion float32 `json:"completion" binding:"required,gte=0,lte=100"`
}

// Finds Todo object from database for given Id. Replaces its parameters
// with those from request and stores updated object back to database.
// Completion value must be greater than that present in the database.
func (s *Server) updateTodoCompletionInfo(ctx *gin.Context) {
	req := UpdateTodoCompletionRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	todo, err := s.Models.Todo.GetOneById(req.Id)
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

	res, err := s.Models.Todo.UpdateOne(db.Todo{
		Id:          req.Id,
		Completion:  req.Completion,
		Title:       todo.Title,
		Description: todo.Description,
		IsDone:      todo.IsDone,
		Expiry:      todo.Expiry,
	})
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Updated todo's completion progress",
		Data:    res,
	})
}

// Required object for updateTodoDoneInfo request. Another json fields are ignored.
// Id must be greater then 1. IsDone needs to be true.
type UpdateTodoDoneRequest struct {
	Id     int64 `json:"id" binding:"required,min=1"`
	IsDone bool  `json:"is_done" binding:"required"`
}

// Finds Todo object from database for given Id. Replaces its parameters
// with those from request and stores updated object back to database.
// IsDone field must change object's value from false to true or it throws BadRequest.
func (s *Server) updateTodoDoneInfo(ctx *gin.Context) {
	req := UpdateTodoDoneRequest{}
	fmt.Println("XD")
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	todo, err := s.Models.Todo.GetOneById(req.Id)
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
	res, err := s.Models.Todo.UpdateOne(db.Todo{
		Id:          req.Id,
		IsDone:      req.IsDone,
		Title:       todo.Title,
		Description: todo.Description,
		Completion:  todo.Completion,
		Expiry:      todo.Expiry,
	})
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Message: "Updated todo's status",
		Data:    res,
	})
}

// Request object with uri to be validated. Id must be greater then 0.
type DeleteTodoRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

// Deletes Todo with given Id. Throws NotFound when it affected 0 rows.
func (s *Server) deleteTodo(ctx *gin.Context) {
	req := DeleteTodoRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.Models.Todo.DeleteOne(req.Id)
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

// Request object with queries. Can be omitted.
type GetTodosRequest struct {
	Period string `form:"period"`
}

// Gets slice of Todo objects depending on given uri. If uri field is omitted, it returns
// all todos from database regardless their status. If uri is in [today, tomorrow, week],
// it returns slice of unfinished Todos that expire within a specified period of time.
// Otherwise it throws BadRequest.
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
		startDate := strings.Join(strings.Split(time.Now().Format("2006-01-02"), "-"), "")
		endDate := strings.Join(strings.Split(time.Now().AddDate(0, 0, 1).Format("2006-01-02"), "-"), "")

		todos, err = s.Models.Todo.GetMany(startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		message = "Got all todos for today"

	case "tomorrow":
		startDate := strings.Join(strings.Split(time.Now().AddDate(0, 0, 1).Format("2006-01-02"), "-"), "")
		endDate := strings.Join(strings.Split(time.Now().AddDate(0, 0, 2).Format("2006-01-02"), "-"), "")

		todos, err = s.Models.Todo.GetMany(startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		message = "Got all todos for tomorrow"

	case "week":
		addDays := 8 - int(time.Now().Weekday())
		startDate := strings.Join(strings.Split(time.Now().Format("2006-01-02"), "-"), "")
		endDate := strings.Join(strings.Split(time.Now().AddDate(0, 0, addDays).Format("2006-01-02"), "-"), "")

		todos, err = s.Models.Todo.GetMany(startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		message = "Got all todos for this week"

	case "":
		todos, err = s.Models.Todo.GetAll()
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
