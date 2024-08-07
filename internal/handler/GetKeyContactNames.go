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

// ByName implements sort.Interface for []NameResult based on the Name field.
type ByName []NameResult

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Data.Name < a[j].Data.Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func GetKeyContactName(w http.ResponseWriter, r *http.Request) {
	imageDir := "images"
	defaultImage := "Default.jpg"
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

		// Convert ID to string and set imageName
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
