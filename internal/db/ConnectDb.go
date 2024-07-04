package db

import (
	"database/sql"
	"log"
	"net/http"
)

func InitDb(w http.ResponseWriter, r *http.Request) (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=Smmarp31461013 dbname=mulvansham sslmode=disable"

	// Open connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Error opening database connection", http.StatusInternalServerError)
		log.Fatalf("Error opening database connection: %v", err)
		return nil, err
	}
	return db, nil
}
