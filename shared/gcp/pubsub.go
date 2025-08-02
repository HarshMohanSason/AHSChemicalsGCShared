package gcp

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/pubsub"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
)

var (
	PubSubClient           *pubsub.Client
	PUBSUB_TOPIC_ID        string
	PUBSUB_SUBSCRIPTION_ID string
)

func InitPubSubDebug(ctx context.Context) {
	shared.InitGCPOnce.Do(func() {
		projectID := os.Getenv("GCP_PROJECT_ID")
		if projectID == "" {
			log.Fatalf("GCP_PROJECT_ID env variable not set")
		}
		PUBSUB_TOPIC_ID = os.Getenv("PUBSUB_TOPIC_ID")
		PUBSUB_SUBSCRIPTION_ID = os.Getenv("PUBSUB_SUBSCRIPTION_ID")
		if projectID == "" || PUBSUB_TOPIC_ID == "" || PUBSUB_SUBSCRIPTION_ID == "" {
			log.Fatalf("PubSub env variables not set")
		}
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create PubSub client: %v", err)
		}
		PubSubClient = client
	})
}

func InitPubSubProd(ctx context.Context) {
	shared.InitGCPOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil { 
			log.Fatalf("No project id found for the GCP project: %v", err)
		}
		PUBSUB_TOPIC_ID = LoadSecretsHelper(projectID, "PUBSUB_TOPIC_ID")
		PUBSUB_SUBSCRIPTION_ID = LoadSecretsHelper(projectID, "PUBSUB_SUBSCRIPTION_ID")
		if PUBSUB_TOPIC_ID == "" || PUBSUB_SUBSCRIPTION_ID == "" {
			log.Fatalf("PubSub env variables not set")
		}
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create PubSub client: %v", err)
		}
		PubSubClient = client
	})
}

func PublishMessage(ctx context.Context, data []byte) error {
	if PubSubClient == nil {
		return fmt.Errorf("PubSub client not initialized")
	}

	topic := PubSubClient.Topic(PUBSUB_TOPIC_ID)
	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	_, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func SubscribeToTopic(ctx context.Context, handler func(ctx context.Context, m *pubsub.Message)) error {
	sub := PubSubClient.Subscription(PUBSUB_SUBSCRIPTION_ID)
	err := sub.Receive(ctx, handler)
	if err != nil {
		return err
	}
	return nil
}
