package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

type Updates struct {
	Message   string `json:"Message"`
	Timestamp string `json:"Timestamp"`
}

func GetUpdates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	currentTime := vars["currenttime"]

	// Parse the currentTime to time.Time
	currentDate, err := time.Parse("2006-01-02", currentTime)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		log.Printf("Invalid date format: %v", err)
		return
	}
	currentDate = currentDate.AddDate(0, 0, +1)
	// Calculate the date one week before the currentDate
	oneWeekBefore := currentDate.AddDate(0, 0, -7)

	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Printf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()

	query := `SELECT "message", "time" FROM "Updates" WHERE "time" BETWEEN $1 AND $2`

	// Execute query with date range parameters
	rows, err := db.Query(query, oneWeekBefore, currentDate)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	var results []Updates

	for rows.Next() {
		var result Updates
		var timestamp time.Time

		if err := rows.Scan(&result.Message, &timestamp); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Printf("Error scanning row: %v", err)
			return
		}

		// Format the timestamp to YYYY-MM-DD
		result.Timestamp = timestamp.Format("2006-01-02")
		results = append(results, result)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through rows", http.StatusInternalServerError)
		log.Printf("Error iterating through rows: %v", err)
		return
	}

	if len(results) == 0 {
		noUpdatesMessage := map[string]string{"message": "No updates"}
		jsonData, err := json.Marshal(noUpdatesMessage)
		if err != nil {
			http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
			log.Printf("Error marshalling JSON: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}

	// Convert results slice to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	// Set Content-Type header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
