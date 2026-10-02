package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Lesson struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AggregatedDashboard struct {
	Service      string `json:"service"`
	PodName      string `json:"pod_name"`
	Lessons      []Lesson `json:"lessons"`
	UserService  string `json:"user_service_status"`
	MediaService string `json:"media_service_status"`
}

func checkService(url string) string {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url + "/healthz")
	if err != nil {
		return fmt.Sprintf("Unavailable: %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return fmt.Sprintf("Healthy (%s)", string(body))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	userSvcURL := os.Getenv("USER_SERVICE_URL")
	if userSvcURL == "" {
		userSvcURL = "http://user-service"
	}

	mediaSvcURL := os.Getenv("MEDIA_SERVICE_URL")
	if mediaSvcURL == "" {
		mediaSvcURL = "http://media-service"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/api/lessons", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		dashboard := AggregatedDashboard{
			Service: "Lesson Service (Microservice Orchestrator)",
			PodName: hostname,
			Lessons: []Lesson{
				{ID: "les-1", Title: "Kubernetes Zero to One", Description: "Pods, Deployments, and Core Networking"},
				{ID: "les-2", Title: "Multi-Microservice Architecture", Description: "Inter-Service DNS and Scaling"},
			},
			UserService:  checkService(userSvcURL),
			MediaService: checkService(mediaSvcURL),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dashboard)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "Lesson Service running on Pod: %s\n", hostname)
	})

	log.Printf("Lesson Service listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
