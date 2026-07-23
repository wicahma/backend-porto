package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 1. Initialize the built-in router (Go 1.22+)
	mux := http.NewServeMux()

	// 2. Register a basic health check route
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Homelab API is running!")
	})

	// 3. Define the server port
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)

	// 4. Start listening
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}