package util

import (
	"database/sql"
	"log"
	"time"
)

func InsertMessageIntoUpdates(message string, db *sql.DB) bool {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Prepare the SQL statement
	query := `INSERT INTO "Updates" (message, time) VALUES ($1, $2)`
	_, err := db.Exec(query, message, currentTime)
	if err != nil {
		log.Printf("Error inserting into db: %v", err)
		return false
	}

	return true
}
