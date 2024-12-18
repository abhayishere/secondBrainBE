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
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	User        string `json:"user"`
	ID          string `json:"id"`
}

type ScreenshotData struct {
	Screenshot string `json:"screenshot"`
	User       string `json:"user"`
	ID         string `json:"id"`
}

func verifyIDToken(idToken string) (string, error) {
	auth, err := app.Auth(context.Background())
	if err != nil {
		fmt.Println("Error creating auth")
		return "", err
	}
	token, err := auth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		fmt.Println("error extracting token")
		return "", err
	}
	return token.UID, nil
}

func saveLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	uid, err := verifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	docRef, _, err := client.Collection("links").Add(ctx, map[string]interface{}{
		"url":         linkData.URL,
		"title":       linkData.Title,
		"description": linkData.Description,
		"user":        uid,
	})
	if err != nil {
		http.Error(w, "Failed to save link data", http.StatusInternalServerError)
		return
	}
	linkData.ID = docRef.ID
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
	authHeader := r.Header.Get("Authorization")
	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	uid, err := verifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	docRef, _, err := client.Collection("screenshots").Add(ctx, map[string]interface{}{
		"screenshot": screenshotData.Screenshot,
		"user":       uid,
	})
	if err != nil {
		http.Error(w, "Failed to save screenshot data", http.StatusInternalServerError)
		return
	}
	screenshotData.ID = docRef.ID
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
	fmt.Println("Reached captain")
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	uid, err := verifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	iter := client.Collection("links").Where("user", "==", uid).Documents(ctx)
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
		link.ID = doc.Ref.ID
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
	authHeader := r.Header.Get("Authorization")
	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	uid, err := verifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	iter := client.Collection("screenshots").Where("user", "==", uid).Documents(ctx)
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
		screenshot.ID = doc.Ref.ID
		screenshots = append(screenshots, screenshot)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(screenshots)
}
func deleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	uid, err := verifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	linkID := r.URL.Query().Get("id")
	if linkID == "" {
		http.Error(w, "Missing link ID", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	doc, err := client.Collection("links").Doc(linkID).Get(ctx)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}
	if doc.Data()["user"] != uid {
		http.Error(w, "Unauthorized to delete this link", http.StatusUnauthorized)
		return
	}

	_, err = client.Collection("links").Doc(linkID).Delete(ctx)
	if err != nil {
		http.Error(w, "Failed to delete link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Link deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
func deleteScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	uid, err := verifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	screenshotID := r.URL.Query().Get("id")
	if screenshotID == "" {
		http.Error(w, "Missing screenshot ID", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	doc, err := client.Collection("screenshots").Doc(screenshotID).Get(ctx)
	if err != nil {
		http.Error(w, "Screenshot not found", http.StatusNotFound)
		return
	}
	if doc.Data()["user"] != uid {
		http.Error(w, "Unauthorized to delete this screenshot", http.StatusUnauthorized)
		return
	}

	_, err = client.Collection("screenshots").Doc(screenshotID).Delete(ctx)
	if err != nil {
		http.Error(w, "Failed to delete screenshot", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Screenshot deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
