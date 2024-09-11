package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	initFirebase()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/savelink", saveLinkHandler)
	mux.HandleFunc("/api/savescreen", screenshotHandler)
	mux.HandleFunc("/get-links", getLinksHandler)
	mux.HandleFunc("/get-screenshots", getScreenshotsHandler)
	frontendurl := os.Getenv("FEURL")
	if frontendurl == "" {
		frontendurl = "8080" // Default to 8080 if PORT is not set
	}
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{frontendurl}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
