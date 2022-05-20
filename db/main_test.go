package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testQueries Model

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"mock",
		"mock",
		"localhost",
		"8888",
		"mock",
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
