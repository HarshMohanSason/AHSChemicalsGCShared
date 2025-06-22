package shared

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/joho/godotenv"
)

var (
	TWILIO_ACCOUNT_SID      string
	TWILIO_AUTH_TOKEN       string
	TWILIO_FROM_PHONE       string
	TWILIO_RECIPENTS_PHONE  string
	SENDGRID_API_KEY        string
	SENDGRID_FROM_MAIL      string
)

/*
	InitThirdPartyDebug loads third-party service configuration for local debug environment.

	Behavior:
	- Loads environment variables from `./keys/.env`.
	- Initializes global variables for:
		- Twilio (ACCOUNT_SID, AUTH_TOKEN, FROM_PHONE, RECIPENTS_PHONE)
		- SendGrid (API_KEY, FROM_MAIL)

	Logs:
	- Fatal error if `.env` cannot be loaded.
*/
func InitThirdPartyDebug() {
	initOnce.Do(func() {
		err := godotenv.Load("./keys/.env")
		if err != nil {
			log.Fatalf("Error loading the env file: %v", err)
		}
		TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")
		TWILIO_AUTH_TOKEN = os.Getenv("TWILIO_AUTH_TOKEN")
		TWILIO_FROM_PHONE = os.Getenv("TWILIO_FROM_PHONE")
		TWILIO_RECIPENTS_PHONE = os.Getenv("TWILIO_RECIPENTS_PHONE")
		SENDGRID_API_KEY = os.Getenv("SENDGRID_API_KEY")
		SENDGRID_FROM_MAIL = os.Getenv("SENDGRID_FROM_MAIL")

		log.Println("Third-party services initialized in DEBUG environment.")
	})
}

/*
	InitThirdPartyProd loads third-party service configuration for production environment
	from Google Cloud Secret Manager.

	Behavior:
	- Uses metadata server to detect GCP project ID.
	- Loads secrets from Secret Manager for:
		- Twilio (ACCOUNT_SID, AUTH_TOKEN, FROM_PHONE, RECIPENTS_PHONE)
		- SendGrid (API_KEY, FROM_MAIL)

	Env variable names in Secret Manager should match **exactly**:
		- TWILIO_ACCOUNT_SID
		- TWILIO_AUTH_TOKEN
		- TWILIO_FROM_PHONE
		- TWILIO_RECIPENTS_PHONE
		- SENDGRID_API_KEY
		- SENDGRID_FROM_MAIL

	Logs:
	- Fatal error if any secret is missing.
	- Logs success after completion.
*/
func InitThirdPartyProd(ctx context.Context) {
	initOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		loadSecret := func(secretName string) string {
			path := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
			secret, err := GetSecretFromGCP(path)
			if err != nil {
				log.Fatalf("Error fetching secret %s: %v", secretName, err)
			}
			return secret
		}

		TWILIO_ACCOUNT_SID = loadSecret("TWILIO_ACCOUNT_SID")
		TWILIO_AUTH_TOKEN = loadSecret("TWILIO_AUTH_TOKEN")
		TWILIO_FROM_PHONE = loadSecret("TWILIO_FROM_PHONE")
		TWILIO_RECIPENTS_PHONE = loadSecret("TWILIO_RECIPENTS_PHONE")
		SENDGRID_API_KEY = loadSecret("SENDGRID_API_KEY")
		SENDGRID_FROM_MAIL = loadSecret("SENDGRID_FROM_MAIL")

		log.Println("Third-party services initialized in PRODUCTION environment.")
	})
}