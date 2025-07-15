package quickbooks

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

var (
	QUICKBOOKS_CLIENT_ID                string
	QUICKBOOKS_CLIENT_SECRET            string
	QUICKBOOKS_AUTH_CALLBACK_URL        string //Redirect URL when first auth url is called and redirects the user to this
	QUICBOOK_AUTH_CALLBACK_REDIRECT_URL string //Redirects user to a successful login page when quickbooks is authenticated
	QUICKBOOKS_API_URL                  string //base url which contains either the debug or production api url.
)

// InitQuickBooksDebug initializes QuickBooks credentials for the **debug** environment.
//
// It loads environment variables from a local .env file
// and sets the global QuickBooks credential variables.
//
// This function is **intended for local development only**
// Environment variables required in the .env file:
//
//	QUICKBOOKS_DEBUG_CLIENT_ID
//	QUICKBOOKS_DEBUG_CLIENT_SECRET
//	QUICKBOOKS_DEBUG_AUTH_CALLBACK_URL
//	QUICKBOOKS_API_URL
//
// Logs:
//
//	Fatal errors if the .env file cannot be loaded or required variables are missing.
//	Success message if successfully initialized
func InitQuickBooksDebug() {
	QUICKBOOKS_CLIENT_ID = os.Getenv("QUICKBOOKS_DEBUG_CLIENT_ID")
	QUICKBOOKS_CLIENT_SECRET = os.Getenv("QUICKBOOKS_DEBUG_CLIENT_SECRET")
	QUICKBOOKS_AUTH_CALLBACK_URL = os.Getenv("QUICKBOOKS_DEBUG_AUTH_CALLBACK_URL")
	QUICBOOK_AUTH_CALLBACK_REDIRECT_URL = os.Getenv("QUICKBOOKS_DEBUG_AUTH_CALLBACK_REDIRECT_URL")
	QUICKBOOKS_API_URL = os.Getenv("QUICKBOOKS_DEBUG_API_URL")

	if QUICKBOOKS_CLIENT_ID == "" || QUICKBOOKS_CLIENT_SECRET == "" || QUICKBOOKS_AUTH_CALLBACK_URL == "" || QUICKBOOKS_API_URL == "" {
		log.Fatalf("Error initializing QuickBooks credentials: missing required environment variables")
	}
	log.Println("Initialized quickbooks credentials in debug...")
}

// InitQuickBooksProd initializes QuickBooks credentials for the **production** environment.
//
// It retrieves credentials from **Google Secret Manager** using the project ID fetched from
// the GCP Metadata server (only available in GCP environments).
//
// This function is **intended for production deployments** and uses `sync.Once` to ensure
// initialization happens only once during the program lifecycle.
//
// Secrets required in Secret Manager (one latest version each):
//   - QUICKBOOKS_CLIENT_ID
//   - QUICKBOOKS_CLIENT_SECRET
//   - QUICKBOOKS_AUTH_CALLBACK_URL
//   - QUICKBOOKS_API_URL
//
// Parameters:
//
//	ctx - context.Context from the calling function, usually the incoming HTTP request context.
//
// Logs:
//
//	Fatal errors if the project ID cannot be determined, or any of the required secrets cannot be fetched.
//	Success message if successfully initialized
func InitQuickBooksProd(ctx context.Context) {
	shared.InitQuickBooksOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		// Load all QuickBooks-related secrets
		QUICKBOOKS_CLIENT_ID = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_CLIENT_ID")
		QUICKBOOKS_CLIENT_SECRET = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_CLIENT_SECRET")
		QUICKBOOKS_AUTH_CALLBACK_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_AUTH_CALLBACK_URL")
		QUICBOOK_AUTH_CALLBACK_REDIRECT_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_AUTH_CALLBACK_REDIRECT_URL")
		QUICKBOOKS_API_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_API_URL")

		if QUICKBOOKS_CLIENT_ID == "" || QUICKBOOKS_CLIENT_SECRET == "" || QUICKBOOKS_AUTH_CALLBACK_URL == "" || QUICKBOOKS_API_URL == "" {
			log.Fatalf("Error initializing QuickBooks credentials: missing required environment variables")
		}
		log.Println("QuickBooks credentials initialized for PRODUCTION environment.")
	})
}
