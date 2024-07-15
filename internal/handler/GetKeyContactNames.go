package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

// ByName implements sort.Interface for []NameResult based on the Name field.
type ByName []NameResult

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Data.Name < a[j].Data.Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func GetKeyContactName(w http.ResponseWriter, r *http.Request) {
	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Fatalf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()

	// Prepare query statement
	query := `SELECT "ID", "Name","Generation" FROM mulvansham WHERE "Relationship" = 1 AND "IsKeyContact" = true;`

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
		if err := rows.Scan(&result.ID, &result.Data.Name, &result.Generation); err != nil {
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

	// Sort results by Name
	sort.Sort(ByName(results))

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
