package quickbooks

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
	// QUICKBOOKS_CLIENT_ID holds the QuickBooks OAuth client ID.
	QUICKBOOKS_CLIENT_ID string

	// QUICKBOOKS_CLIENT_SECRET holds the QuickBooks OAuth client secret.
	QUICKBOOKS_CLIENT_SECRET string

	// QUICKBOOKS_AUTH_CALLBACK_URL is the redirect URI used for QuickBooks OAuth callbacks.
	QUICKBOOKS_AUTH_CALLBACK_URL string

	// QUICBOOK_AUTH_CALLBACK_REDIRECT_URL is the redirect URI used to show user a confirmation page when quickbooks is authenticated 
	QUICBOOK_AUTH_CALLBACK_REDIRECT_URL string
	
	// QUICKBOOKS_API_URL is the base url which contains either the debug or production api url.
	QUICKBOOKS_API_URL string

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
// Logs:
//
//	Fatal errors if the .env file cannot be loaded or required variables are missing.
//	Success message if successfully initialized
//
// Example usage (local development):
//
//	shared.InitQuickBooksDebug()
func InitQuickBooksDebug() {
	QUICKBOOKS_CLIENT_ID = os.Getenv("QUICKBOOKS_DEBUG_CLIENT_ID")
	QUICKBOOKS_CLIENT_SECRET = os.Getenv("QUICKBOOKS_DEBUG_CLIENT_SECRET")
	QUICKBOOKS_AUTH_CALLBACK_URL = os.Getenv("QUICKBOOKS_DEBUG_AUTH_CALLBACK_URL")
	QUICBOOK_AUTH_CALLBACK_REDIRECT_URL = os.Getenv("QUICKBOOKS_DEBUG_AUTH_CALLBACK_REDIRECT_URL")
	QUICKBOOKS_API_URL = os.Getenv("QUICKBOOKS_DEBUG_API_URL")	
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
//
// Example usage (in production):
//
//	shared.InitQuickBooksProd(ctx)
func InitQuickBooksProd(ctx context.Context) {
	shared.InitQuickBooksOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		// Helper function to load a secret by its name
		loadSecret := func(secretName string) string {
			path := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
			secret, err := gcp.GetSecretFromGCP(path)
			if err != nil {
				log.Fatalf("Error fetching secret %s: %v", secretName, err)
			}
			return secret
		}

		// Load all QuickBooks-related secrets
		QUICKBOOKS_CLIENT_ID = loadSecret("QUICKBOOKS_CLIENT_ID")
		QUICKBOOKS_CLIENT_SECRET = loadSecret("QUICKBOOKS_CLIENT_SECRET")
		QUICKBOOKS_AUTH_CALLBACK_URL = loadSecret("QUICKBOOKS_AUTH_CALLBACK_URL")
		QUICBOOK_AUTH_CALLBACK_REDIRECT_URL = loadSecret("QUICKBOOKS_AUTH_CALLBACK_REDIRECT_URL")
		QUICKBOOKS_API_URL = loadSecret("QUICKBOOKS_API_URL")
		log.Println("QuickBooks credentials initialized for PRODUCTION environment.")
	})
}
