package main

import (
	"finance-crud-app/internal/db"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:Password123@localhost:5432/crud_db?sslmode=disable"

	dbconn, err := db.NewPGStorage(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	db.SeedTestDB(dbconn)
}
