package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Item struct {
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Available int    `json:"available"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/api/inventory", func(w http.ResponseWriter, r *http.Request) {
		items := []Item{
			{SKU: "K8S-101", Name: "Kubernetes Zero to Hero Course", Available: 50},
			{SKU: "GO-201", Name: "Advanced Go Microservices", Available: 25},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "Inventory Service running on Pod: %s\n", hostname)
	})

	log.Printf("Inventory Service listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
