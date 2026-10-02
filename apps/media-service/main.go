package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type MediaAsset struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
	URL      string `json:"url"`
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

	http.HandleFunc("/api/media", func(w http.ResponseWriter, r *http.Request) {
		media := []MediaAsset{
			{ID: "m-1", Title: "Kubernetes Architecture Video", Duration: "15m30s", URL: "https://stream.example.com/k8s-arch.mp4"},
			{ID: "m-2", Title: "Pods and Services Deep Dive", Duration: "22m10s", URL: "https://stream.example.com/k8s-pods.mp4"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(media)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "Media Service running on Pod: %s\n", hostname)
	})

	log.Printf("Media Service listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
