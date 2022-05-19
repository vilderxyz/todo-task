package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vilderxyz/todos/db"
	"gorm.io/gorm"
)

type Config struct {
	ServerAddr string
}

// Main struct of the application that handles traffic
// and is responsible of database connection
type Server struct {
	Models db.Models
	Router *gin.Engine
}

// Creates a new Server instance with database connection
// and returns pointer to it
func NewServer(conn *gorm.DB) *Server {
	server := &Server{
		Models: db.New(conn),
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
// Prints error and stores it
// in the map returned to the client's side
func errorResponse(err error) gin.H {
	log.Println(err)
	return gin.H{"error": err.Error()}
}
