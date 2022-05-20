package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/vilderxyz/todos/db"
	valid "github.com/vilderxyz/todos/validator"
	"gorm.io/gorm"
)

// Struct of http server for Todos application.
type Server struct {
	Queries db.DB
	Router  *gin.Engine
}

// Creates a new Server instance with database connection
// and returns pointer to it
func NewServer(conn *gorm.DB) *Server {
	server := &Server{
		Queries: db.New(conn),
	}

	// Registers custom period validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("period", valid.ValidPeriod)
	}

	server.setupRouter()
	return server
}

// Runs Gin router on given address
func (s *Server) Start(addr string) error {
	log.Println("Serving at: ", addr)
	return s.Router.Run(addr)
}

// Helps handling errors much faster.
//
// Prints an error and sends it back to the client's side
func errorResponse(err error) gin.H {
	log.Println(err)
	return gin.H{"error": err.Error()}
}
