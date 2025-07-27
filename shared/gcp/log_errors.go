package gcp

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
)

var (
	once        sync.Once
	LogClient   *logging.Client
)

func InitLogger(ctx context.Context) {
	once.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Failed to get project ID: %v", err)
		}

		LogClient, err = logging.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create logger: %v", err)
		}
	})
}

func CloseLogger() {
	if LogClient != nil {
		_ = LogClient.Close()
	}
}