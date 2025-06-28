package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

var (
	TWILIO_ACCOUNT_SID      string
	TWILIO_AUTH_TOKEN       string
	TWILIO_FROM_PHONE       string
	TWILIO_RECIPIENTS_PHONE string
	SENDGRID_API_KEY        string
	SENDGRID_FROM_MAIL      string
)

/*
InitThirdPartyDebug initializes third-party service configuration from local environment variables.

It assumes a `.env` file is already loaded (e.g., via dotenv or shell).

Expected ENV Variables:
- TWILIO_ACCOUNT_SID
- TWILIO_AUTH_TOKEN
- TWILIO_FROM_PHONE
- TWILIO_RECIPIENTS_PHONE
- SENDGRID_API_KEY
- SENDGRID_FROM_MAIL

Logs:
- Logs success or fatal error if required variables are missing.
*/
func InitThirdPartyDebug() {
	TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_AUTH_TOKEN = os.Getenv("TWILIO_AUTH_TOKEN")
	TWILIO_FROM_PHONE = os.Getenv("TWILIO_FROM_PHONE")
	TWILIO_RECIPIENTS_PHONE = os.Getenv("TWILIO_RECIPIENTS_PHONE")
	SENDGRID_API_KEY = os.Getenv("SENDGRID_API_KEY")
	SENDGRID_FROM_MAIL = os.Getenv("SENDGRID_FROM_MAIL")

	// Optional: validate presence
	if TWILIO_ACCOUNT_SID == "" || SENDGRID_API_KEY == "" {
		log.Fatal("One or more required third-party environment variables are missing in DEBUG mode")
	}

	log.Println("Third-party services initialized in DEBUG environment.")
}

/*
InitThirdPartyProd loads third-party service configuration from Google Cloud Secret Manager.

It auto-detects GCP project ID using metadata server and reads the following secrets:

Required Secrets in Secret Manager:
- TWILIO_ACCOUNT_SID
- TWILIO_AUTH_TOKEN
- TWILIO_FROM_PHONE
- TWILIO_RECIPIENTS_PHONE
- SENDGRID_API_KEY
- SENDGRID_FROM_MAIL

Logs:
- Fatal error if any secret is missing.
- Logs successful initialization.
*/
func InitThirdPartyProd(ctx context.Context) {
	shared.InitThirdPartyOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Failed to get GCP project ID: %v", err)
		}

		loadSecret := func(name string) string {
			secretPath := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, name)
			secret, err := gcp.GetSecretFromGCP(secretPath)
			if err != nil {
				log.Fatalf("❌ Error fetching secret %s: %v", name, err)
			}
			return secret
		}

		TWILIO_ACCOUNT_SID = loadSecret("TWILIO_ACCOUNT_SID")
		TWILIO_AUTH_TOKEN = loadSecret("TWILIO_AUTH_TOKEN")
		TWILIO_FROM_PHONE = loadSecret("TWILIO_FROM_PHONE")
		TWILIO_RECIPIENTS_PHONE = loadSecret("TWILIO_RECIPIENTS_PHONE")
		SENDGRID_API_KEY = loadSecret("SENDGRID_API_KEY")
		SENDGRID_FROM_MAIL = loadSecret("SENDGRID_FROM_MAIL")

		log.Println("Third-party services initialized in PRODUCTION environment.")
	})
}
