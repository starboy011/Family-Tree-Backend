package handler

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

type Gallery struct {
	ID   int
	Name string
	Img  string
}

func GetImageGallery(w http.ResponseWriter, r *http.Request) {
	imageDir := "images/ImageGallery"
	defaultImage := "Default.jpg"
	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Fatalf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()
	query := `SELECT "id", "name" FROM "ImageGallery"`

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		log.Fatalf("Error executing query: %v", err)
		return
	}
	defer rows.Close()
	var results []Gallery

	for rows.Next() {
		var result Gallery
		if err := rows.Scan(&result.ID, &result.Name); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Fatalf("Error scanning row: %v", err)
			return
		}
		idStr := strconv.Itoa(result.ID)
		imageName := idStr + ".jpg"
		imagePath := path.Join(imageDir, imageName)

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
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		result.Img = "data:image/jpg;base64," + imageBase64
		sort.Slice(results, func(i, j int) bool {
			return results[i].ID < results[j].ID
		})

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
