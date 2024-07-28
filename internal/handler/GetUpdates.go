package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

type Updates struct {
	Message   string `json:"Message"`
	Timestamp string `json:"Timestamp"`
}

func GetUpdates(w http.ResponseWriter, r *http.Request) {
	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Fatalf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()
	query := `SELECT "message", "time" FROM "Updates"`

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		log.Fatalf("Error executing query: %v", err)
		return
	}
	defer rows.Close()
	var results []Updates

	for rows.Next() {
		var result Updates
		var timestamp time.Time

		if err := rows.Scan(&result.Message, &timestamp); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatalf("Error scanning row: %v", err)
			return
		}

		// Format the timestamp to YYYY-MM-DD
		result.Timestamp = timestamp.Format("2006-01-02")
		results = append(results, result)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through rows", http.StatusInternalServerError)
		log.Fatalf("Error iterating through rows: %v", err)
		return
	}

	// Convert results slice to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		log.Fatalf("Error marshalling JSON: %v", err)
		return
	}

	// Set Content-Type header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
