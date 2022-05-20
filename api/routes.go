package api

import (
	"github.com/gin-gonic/gin"
)

// Setups all available http routes
func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/todos", s.createTodo)

	router.GET("/todos/:id", s.getTodoById)

	router.GET("/todos", s.getTodos)

	router.PATCH("/todos", s.updateTodoTextInfo)

	router.PATCH("/todos/completion", s.updateTodoCompletionInfo)

	router.PATCH("/todos/done", s.updateTodoDoneInfo)

	router.DELETE("/todos/:id", s.deleteTodo)

	s.Router = router
}
