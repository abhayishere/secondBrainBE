package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	initFirebase()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/savelink", saveLinkHandler)
	mux.HandleFunc("/api/savescreen", screenshotHandler)
	mux.HandleFunc("/get-links", getLinksHandler)
	mux.HandleFunc("/get-screenshots", getScreenshotsHandler)
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Allow your frontend URL
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(mux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
