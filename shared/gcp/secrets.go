package gcp

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
)

var (
	secretManagerClient *secretmanager.Client
	secretManagerErr    error
)

// getSecretManagerClient initializes the Secret Manager client once.
func getSecretManagerClient(ctx context.Context) (*secretmanager.Client, error) {
	shared.InitGCPOnce.Do(func() {
		secretManagerClient, secretManagerErr = secretmanager.NewClient(ctx)
	})
	return secretManagerClient, secretManagerErr
}

// GetSecretFromGCP retrieves the latest version of a secret from Google Cloud Secret Manager.
//
// Parameters:
//   - secretName: The fully-qualified name of the secret version in the format:
//     "projects/{projectID}/secrets/{secretName}/versions/{version}"
//     For example: "projects/my-project/secrets/my-secret/versions/latest"
//
// Returns:
//   - The secret payload as a string if retrieval is successful.
//   - An error if any occurs during client creation or secret access.
func GetSecretFromGCP(secretName string) (string, error) {
	ctx := context.Background()

	// Initialize the Secret Manager client.
	client, err := getSecretManagerClient(ctx)
	if err != nil {
		return "", err
	}

	// Build the request to access the secret version.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	// Access the secret version and extract its payload.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	// Convert the payload from bytes to string and return.
	secretData := string(result.Payload.Data)
	return secretData, nil
}
