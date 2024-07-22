package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/starboy011/Family-Tree-Backend/internal/db"
)

func PushFcmToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Printf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO fcmtoken (token) VALUES ($1)", token)
	if err != nil {
		// Check for duplicate key error
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			// Handle duplicate key error
			response := map[string]string{"message": "Token already exists"}
			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Error in creating response", http.StatusInternalServerError)
				log.Printf("Error creating response: %v", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
			return
		}

		http.Error(w, "Error inserting token into db", http.StatusInternalServerError)
		log.Printf("Error inserting token into db: %v", err)
		return
	}

	// Respond with success
	response := map[string]string{"message": "Token successfully saved"}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error in creating response", http.StatusInternalServerError)
		log.Printf("Error creating response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
