package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/api/iterator"
)

type LinkData struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type ScreenshotData struct {
	Screenshot string `json:"screenshot"`
}

func saveLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var linkData LinkData
	err = json.Unmarshal(body, &linkData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Save to Firestore
	ctx := context.Background()
	_, _, err = client.Collection("links").Add(ctx, map[string]interface{}{
		"url":   linkData.URL,
		"title": linkData.Title,
	})
	if err != nil {
		http.Error(w, "Failed to save link data", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received URL: %s, Title: %s\n", linkData.URL, linkData.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Link saved successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func screenshotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var screenshotData ScreenshotData
	err = json.Unmarshal(body, &screenshotData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	base64Image := screenshotData.Screenshot[strings.IndexByte(screenshotData.Screenshot, ',')+1:]
	imgData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusBadRequest)
		return
	}
	imgFilePath := "screenshot.png"
	err = ioutil.WriteFile(imgFilePath, imgData, 0644)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Save to Firestore
	ctx := context.Background()
	_, _, err = client.Collection("screenshots").Add(ctx, map[string]interface{}{
		"screenshot": screenshotData.Screenshot,
	})
	if err != nil {
		http.Error(w, "Failed to save screenshot data", http.StatusInternalServerError)
		return
	}

	fmt.Println("Screenshot received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Screenshot saved successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func getLinksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	iter := client.Collection("links").Documents(ctx)
	defer iter.Stop()

	var links []LinkData
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to retrieve links", http.StatusInternalServerError)
			return
		}

		var link LinkData
		doc.DataTo(&link)
		links = append(links, link)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
func getScreenshotsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	iter := client.Collection("screenshots").Documents(ctx)
	defer iter.Stop()

	var screenshots []ScreenshotData
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to retrieve screenshots", http.StatusInternalServerError)
			return
		}

		var screenshot ScreenshotData
		doc.DataTo(&screenshot)
		screenshots = append(screenshots, screenshot)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(screenshots)
}
