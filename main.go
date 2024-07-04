package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/starboy011/Family-Tree-Backend/internal/handler"
)

func main() {
	// Define HTTP route for /get-all-names endpoint
	http.HandleFunc("/get-all-names", handler.GetAllName)

	// Start HTTP server on port 8080
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
