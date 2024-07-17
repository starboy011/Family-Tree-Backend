package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/starboy011/Family-Tree-Backend/internal/db"
	"github.com/starboy011/Family-Tree-Backend/pkg"
)

func GetIdFromName(w http.ResponseWriter, r *http.Request, firstName string) int {
	// Initialize database
	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		log.Printf("Error connecting to db: %v", err)
		return 0
	}
	defer db.Close()

	// Prepare query statement
	query := `SELECT "ID" FROM mulvansham WHERE "Relationship" = 1 AND "Name" = $1;`

	// Execute query
	var id int
	err = db.QueryRow(query, firstName).Scan(&id)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		log.Printf("Error executing query: %v", err)
		return 0
	}

	return id
}

func GetRelationshipFromName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	firstName := vars["firstname"]
	secondName := vars["secondname"]

	firstId := GetIdFromName(w, r, firstName)
	secondId := GetIdFromName(w, r, secondName)
	firstTreeData := pkg.GetTreeDataFromId(w, r, firstId)
	secondTreeData := pkg.GetTreeDataFromId(w, r, secondId)
	fmt.Println(firstTreeData, secondTreeData)
}
