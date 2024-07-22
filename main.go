package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/starboy011/Family-Tree-Backend/internal/handler"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define HTTP routes
	r.HandleFunc("/names", handler.GetAllName).Methods("GET")
	r.HandleFunc("/keyNames", handler.GetKeyContactName).Methods("GET")
	r.HandleFunc("/tree", handler.GetFamilyTreeData).Methods("GET")
	r.HandleFunc("/tree/{id}", handler.GetTreeDataFromId).Methods("GET")
	r.HandleFunc("/generation/{id}", handler.GetListFromGeneration).Methods("GET")
	r.HandleFunc("/relationship/{firstname}/{secondname}", handler.GetRelationshipFromName).Methods("GET")
	r.HandleFunc("/fcm/{token}", handler.PushFcmToken).Methods("GET")
	// Start HTTP server on port 8080
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
