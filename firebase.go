package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var app *firebase.App
var client *firestore.Client

func initFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("/Users/abhayishere/Downloads/SecondBrain/backend/secondbrain-6ba99-firebase-adminsdk-221bw-3bbdda7e16.json")

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
