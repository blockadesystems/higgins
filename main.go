package main

import (
	"log"
	"net/http"
)

func main() {
	// Start the Raft node
	node := NewRaftNode()

	// Create the API
	api := NewRaftAPI(node)

	// Set up the HTTP routes
	http.HandleFunc("/get", api.GetValue)
	http.HandleFunc("/set", api.SetValue)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
