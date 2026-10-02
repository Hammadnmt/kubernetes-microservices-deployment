package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Health check endpoint for Kubernetes
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 2. User API endpoint
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		users := []User{
			{ID: "1", Name: "Hammad", Email: "hammad@example.com", Role: "Admin"},
			{ID: "2", Name: "Alice", Email: "alice@example.com", Role: "Developer"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	// 3. Root greeting
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "User Service is running on Pod: %s\n", hostname)
	})

	log.Printf("User Service listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
