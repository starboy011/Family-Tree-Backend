package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// NameResult represents the structure of data to return as JSON
type NameResult struct {
	ID   int `json:"id"`
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

// GetAllName fetches ID and Name columns from the mulvansham table and returns as JSON
func GetAllName(w http.ResponseWriter, r *http.Request) {
	// PostgreSQL connection string
	connStr := "host=localhost port=5432 user=postgres password=Smmarp31461013 dbname=mulvansham sslmode=disable"

	// Open connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Error opening database connection", http.StatusInternalServerError)
		log.Fatalf("Error opening database connection: %v", err)
		return
	}
	defer db.Close()

	// Prepare query statement
	query := `SELECT "ID", "Name" FROM mulvansham`

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		log.Fatalf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	// Create a slice to hold the results
	var results []NameResult

	// Iterate through rows and populate results
	for rows.Next() {
		var result NameResult
		if err := rows.Scan(&result.ID, &result.Data.Name); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatalf("Error scanning row: %v", err)
			return
		}
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
