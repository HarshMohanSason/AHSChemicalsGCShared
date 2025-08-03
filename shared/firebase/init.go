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

// Firebase global variables initialized only once at runtime
var (
	App              *firebase.App
	AuthClient       *auth.Client
	StorageClient    *storage.Client
	FirestoreClient  *firestore.Client
	StorageBucket    string
	initFirebaseOnce sync.Once
)

// InitFirebaseDebug initializes Firebase clients for the debug environment.
//
// Parameters:
//   - keyPath: Path to the Firebase Admin SDK service account JSON file.
//
// Behavior:
//   - Initializes the Firebase App with the specified debug Realtime Database URL.
//   - Sets up AuthClient, FirestoreClient, RealtimeClient, and StorageClient for subsequent usage.
//   - Uses sync.Once to ensure initialization happens only once during the application lifecycle.
//
// Logs:
//   - Calls log.Fatalf() and exits the application if initialization of any service fails.
func InitFirebaseDebug(keyPath string) {
	initFirebaseOnce.Do(func() {
		ctx := context.Background()
		var err error

		// Initialize Firebase App with credentials file.
		opt := option.WithCredentialsFile(keyPath)
		App, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			log.Fatalf("Error occurred initializing Firebase: %v", err)
		}

		// Initialize individual Firebase service clients.
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

// InitFirebaseStaging initializes Firebase clients for the staging environment.
//
// Parameters:
//   - keyPath: Optional pointer to the path of the Firebase Admin SDK service account JSON file.
//     If nil, the default credentials will be used.
//
// Behavior:
//   - Sets up AuthClient, FirestoreClient, and StorageClient for subsequent usage.
//   - Uses sync.Once to ensure initialization happens only once during the application lifecycle.
//
// Logs:
//   - Calls log.Fatalf() and exits the application if initialization of any service fails.
func InitFirebaseStaging(keyPath *string) {
	initFirebaseOnce.Do(func() {
		ctx := context.Background()
		var err error

		// Initialize Firebase App with provided credentials or fallback to default.
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

		// Initialize individual Firebase service clients.
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
		//Get the storage bucket url from gcp secret manager
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

// InitFirebaseProd initializes Firebase clients for the production environment.
//
// Parameters:
//   - keyPath: Optional pointer to the path of the Firebase Admin SDK service account JSON file.
//     If nil, the default credentials will be used.
//
// Behavior:
//   - Sets up AuthClient, FirestoreClient, and StorageClient for subsequent usage.
//   - Uses sync.Once to ensure initialization happens only once during the application lifecycle.
//
// Logs:
//   - Calls log.Fatalf() and exits the application if initialization of any service fails.
func InitFirebaseProd(keyPath *string) {
	initFirebaseOnce.Do(func() {
		ctx := context.Background()
		var err error

		// Initialize Firebase App with provided credentials or fallback to default.
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

		// Initialize individual Firebase service clients.
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
		//Get the storage bucket url from gcp secret manager
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
