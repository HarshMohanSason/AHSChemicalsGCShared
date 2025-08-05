package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/pubsub"
)

var (
	PubSubClient           *pubsub.Client
	PUBSUB_TOPIC_ID        string
	PUBSUB_SUBSCRIPTION_ID string
	initPubSubOnce         sync.Once
)

func InitPubSubDebug(ctx context.Context) {
	initPubSubOnce.Do(func() {
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
		log.Println("PubSub initialized")
	})
}

func InitPubSubFromSecrets(ctx context.Context) {
	initPubSubOnce.Do(func() {
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
		log.Println("PubSub initialized")
	})
}

// Represents the data payload received by a subscriber to a topic. The `Data` bytes
// can then be accessed to decode it to the struct that the `Data` was sent initially 
// as marshalled.
type SubMessage struct {
	Message struct {
		Data        []byte            `json:"data"`        
		MessageID   string            `json:"messageId"`
		PublishTime string            `json:"publishTime"`
		Attributes  map[string]string `json:"attributes"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// DecodeSubMessageData decodes the message data as a T struct. Returns a pointer of
// type T struct. Using T for future usecases.
func DecodeSubMessageData[T any](msg *SubMessage) (*T, error) {
	var result T
	if err := json.Unmarshal(msg.Message.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

//PublishMessage publishes a message to a PubSub topic. Returns an error if
//failed to publish the message. PubSubClient must be initialized before calling
//this function.
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
