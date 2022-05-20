package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/vilderxyz/todos/mock"
)

func newTestServer(t *testing.T, mockModel *mock.MockDB) *Server {
	server := NewServer(nil)
	require.NotEmpty(t, server)

	server.Queries = mockModel
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
