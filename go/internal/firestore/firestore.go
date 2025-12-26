package firestore

// like mentioned in auth.go, we are committing to using goth and gothic session management
// we will be CRUD server side instead of client side fetch/writes
// which are implemented here in this file

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

var Client *firestore.Client

// NOTE: DO NOT use json to initialize the client unless you are ready to pay for secret manager
// (it is okay to do in local dev, but bad practice for production and you will have to pay for it)
//
// you will suffer and want to delete everything when you are deploying to GCP under free tier
// hours spent realizing this mistake: ~8

// but in google cloud functions / app engine, the client initializes automatically
// with the service account attached to the instance
// so we will use that method of initialization for production

// init firestore client
func Init() error {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	if projectID == "" {
		return errors.New("project ID is not set in env variables")
	}

	var err error
	// using background context so the client lives as long as the app lives
	Client, err = firestore.NewClient(context.Background(), projectID)
	if err != nil {
		return err
	}

	log.Println("firestore initialized successfully")
	return nil
}

// clean up the firestore client
func Close() {
	if Client != nil {
		Client.Close()
	}
}
