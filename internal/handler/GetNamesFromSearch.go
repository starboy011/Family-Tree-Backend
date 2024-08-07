package handler

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

func GetName(w http.ResponseWriter, r *http.Request) {
	imageDir := "images"
	defaultImage := "Default.jpg"
	vars := mux.Vars(r)
	Name := vars["name"]

	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `SELECT "ID", "Name", "Generation" FROM mulvansham WHERE "Name" LIKE '%' || $1 || '%' AND "Relationship" <> 3`

	// Execute query with the name parameter
	rows, err := db.Query(query, Name)
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

		// Read image file as bytes
		idStr := strconv.Itoa(result.ID)
		imageName := idStr + ".jpg"
		imagePath := path.Join(imageDir, imageName)

		// Check if the image file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			// Use default image if the image file does not exist
			imagePath = path.Join(imageDir, defaultImage)
		}

		imageBytes, err := os.ReadFile(imagePath)
		if err != nil {
			http.Error(w, "Error reading image file", http.StatusInternalServerError)
			log.Fatalf("Error reading image file: %v", err)
			return
		}

		// Encode image bytes to base64
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		result.Data.Img = "data:image/jpg;base64," + imageBase64 // Adjust according to your image type

		results = append(results, result)
	}

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

	// Set Content-Type header for JSON response
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}
