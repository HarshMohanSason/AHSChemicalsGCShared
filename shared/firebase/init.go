package firebase_shared

import (
	"context"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/storage"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
	"google.golang.org/api/option"
)

var (
	App              *firebase.App     //Firestore App
	AuthClient       *auth.Client      //Auth Client
	StorageClient    *storage.Client   //Storage Client
	FirestoreClient  *firestore.Client //Firestore Client
	StorageBucket    string            //Storage Bucket url
	initFirebaseOnce sync.Once
)

func InitFirebaseDebug(keyPath string) {
	initFirebaseOnce.Do(func() {
		if keyPath == "" {
			log.Fatalf("KeyPath is empty. Required for initializing firebase in debug")
		}
		ctx := context.Background()
		var err error

		opt := option.WithCredentialsFile(keyPath)
		App, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			log.Fatalf("Error occurred initializing Firebase: %v", err)
		}

		AuthClient, err = App.Auth(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Auth client: %v", err)
		}

		FirestoreClient, err = App.Firestore(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Firestore client: %v", err)
		}

		StorageClient, err = App.Storage(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Storage client: %v", err)
		}
		StorageBucket = os.Getenv("STORAGE_BUCKET")
		if StorageBucket == "" {
			log.Fatalf("STORAGE_BUCKET environment variable not set")
		}
	})
}

//Initializes firebase production and staging. Key path is optional since google cloud functions 
//automatically initialize to the correct firebase project they have. When testing in staging 
//environments, `admin_sdk_json` can be passed to do operations like giving admin privileges, 
//db migrations etc.
func InitFirebaseFromSecrets(keyPath *string) {
	initFirebaseOnce.Do(func() {
		ctx := context.Background()
		var err error

		var app *firebase.App
		if keyPath != nil {
			opt := option.WithCredentialsFile(*keyPath)
			app, err = firebase.NewApp(ctx, nil, opt)
		} else {
			app, err = firebase.NewApp(ctx, nil)
		}
		if err != nil {
			log.Fatalf("Error initializing Firebase: %v", err)
		}
		App = app

		AuthClient, err = App.Auth(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Auth client: %v", err)
		}

		FirestoreClient, err = App.Firestore(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Firestore client: %v", err)
		}

		StorageClient, err = App.Storage(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Storage client: %v", err)
		}

		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}
		StorageBucket = gcp.LoadSecretsHelper(projectID, "STORAGE_BUCKET")
		if StorageBucket == "" {
			log.Fatalf("Failed to initialize Storage client: %v", err)
		}
	})
}
