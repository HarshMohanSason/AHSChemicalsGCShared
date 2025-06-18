package shared

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

var (
	App             *firebase.App
	AuthClient      *auth.Client
	StorageClient   *storage.Client
	FirestoreClient *firestore.Client
	RealtimeClient  *db.Client
	initOnce        sync.Once
)

func InitFirebaseDebug(keyPath string) {
	initOnce.Do(func() {
		ctx := context.Background()
		var err error

		//Realtime database url
		conf := &firebase.Config{
			DatabaseURL: "https://ahschemicalsdebug-default-rtdb.firebaseio.com/",
		}

		opt := option.WithCredentialsFile(keyPath)
		App, err = firebase.NewApp(ctx, conf, opt)

		if err != nil {
			log.Fatalf("Error occurred intializing firebase: %v", err)
		}

		//Initialize the Auth Client
		AuthClient, err = App.Auth(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Auth client: %v", err)
		}

		//Initialize the Firestore Client
		FirestoreClient, err = App.Firestore(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Firestore client: %v", err)
		}

		//Initialize the RealtimeDb Client
		RealtimeClient, err = App.Database(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Realtime db client: %v", err)
		}

		//Initialize the Storage Client
		StorageClient, err = App.Storage(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Storage client: %v", err)
		}
	})
}

func InitFirebaseProd(keyPath *string) {
	initOnce.Do(func() {
		ctx := context.Background()
		var err error

		//Realtime db url
		conf := &firebase.Config{
			DatabaseURL: "https://ahschemicalsprod-default-rtdb.firebaseio.com",
		}

		//Initialize the Firebase App
		var app *firebase.App
		if keyPath != nil {
			opt := option.WithCredentialsFile(*keyPath)
			app, err = firebase.NewApp(ctx, conf, opt)
		} else {
			// Use default credentials
			app, err = firebase.NewApp(ctx, conf)
		}

		if err != nil {
			log.Fatalf("Error initializing firebase: %v", err)
		}

		App = app

		//Initialize the Auth Client
		AuthClient, err = App.Auth(ctx)
		if err != nil {
			log.Fatalf("Error initializing auth client: %v", err)
		}

		//Initialize the Firestore Client
		FirestoreClient, err = App.Firestore(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Firestore client: %v", err)
		}

		//Initialize the Realtime Cleint
		RealtimeClient, err = App.Database(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize realtime db: %v", err)
		}

		//Initialize the Storage Client
		StorageClient, err = App.Storage(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Storage client: %v", err)
		}
	})
}
