package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var app *firebase.App
var client *firestore.Client

func initFirebase() {
	ctx := context.Background()
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatalf("Error loading .env file: %v", err1)
	}
	credsbase64 := os.Getenv("FIREBASE_CREDENTIALS")
	if credsbase64 == "" {
		log.Fatalf("FIREBASE_CREDENTIALS environment variable not set")
	}
	credsJSON, err2 := base64.StdEncoding.DecodeString(credsbase64)
	if err2 != nil {
		log.Fatalf("Error decoding Firebase credentials: %v", err2)
	}
	creds := []byte(credsJSON)
	opt := option.WithCredentialsJSON(creds)
	var err error
	app, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v", err)
	}
}
