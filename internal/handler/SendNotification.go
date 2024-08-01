package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/starboy011/Family-Tree-Backend/internal/db"
	"github.com/starboy011/Family-Tree-Backend/internal/util"
)

func SendNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	message := vars["message"]

	db, err := db.InitDb(w, r)
	if err != nil {
		http.Error(w, "Error in connecting to db", http.StatusInternalServerError)
		log.Printf("Error connecting to db: %v", err)
		return
	}
	defer db.Close()

	isIserted := util.InsertMessageIntoUpdates(message, db)
	if !isIserted {
		http.Error(w, "Error while inserting", http.StatusInternalServerError)
		log.Printf("Error while inserting: %v", err)
		return
	}
	tokens := util.GetFcmToken(db)
	for _, token := range tokens.Tokens {
		isNotificationSent := util.SendNotificationViaFirebase(token, title, message)
		if !isNotificationSent {
			http.Error(w, "Error while sending notification via firebase", http.StatusInternalServerError)
			log.Printf("Error while sending notification via firebase: %v", err)
			return
		}
	}
	// Send a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Notification sent successfully")
}
