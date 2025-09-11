// package quickbooks contains initialization of quickbook credentials and basic shared functions
package quickbooks

import (
	"context"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

var (
	QUICKBOOKS_CLIENT_ID                string
	QUICKBOOKS_CLIENT_SECRET            string
	QUICKBOOKS_AUTH_CALLBACK_URL        string //Redirect URL. When the first authentication via initial auth url completes, it redirects to this url
	QUICBOOK_AUTH_CALLBACK_REDIRECT_URL string //Redirects user to a successful login page when quickbooks is authenticated
	QUICKBOOKS_API_URL                  string //base url which contains either the debug or production api url.
	QUICKBOOKS_WEBHOOK_VERIFY_TOKEN     string
	QUICKBOOKS_GET_CUSTOMER_URL         string //func url
	QUICKBOOKS_GET_PRODUCT_URL          string //func url
	QUICKBOOKS_CREATE_ESTIMATE_URL      string //func url
	QUICKBOOKS_DELETE_ESTIMATE_URL      string //func url
	initQuickBooksOnce                  sync.Once
)

func InitQuickBooksDebug() {
	initQuickBooksOnce.Do(func() {
		QUICKBOOKS_CLIENT_ID = os.Getenv("QUICKBOOKS_DEBUG_CLIENT_ID")
		QUICKBOOKS_CLIENT_SECRET = os.Getenv("QUICKBOOKS_DEBUG_CLIENT_SECRET")
		QUICKBOOKS_AUTH_CALLBACK_URL = os.Getenv("QUICKBOOKS_DEBUG_AUTH_CALLBACK_URL")
		QUICBOOK_AUTH_CALLBACK_REDIRECT_URL = os.Getenv("QUICKBOOKS_DEBUG_AUTH_CALLBACK_REDIRECT_URL")
		QUICKBOOKS_API_URL = os.Getenv("QUICKBOOKS_DEBUG_API_URL")
		QUICKBOOKS_GET_CUSTOMER_URL = os.Getenv("QUICKBOOKS_DEBUG_GET_CUSTOMER_URL")
		QUICKBOOKS_GET_PRODUCT_URL = os.Getenv("QUICKBOOKS_DEBUG_GET_PRODUCT_URL")
		QUICKBOOKS_CREATE_ESTIMATE_URL = os.Getenv("QUICKBOOKS_DEBUG_CREATE_ESTIMATE_URL")
		QUICKBOOKS_DELETE_ESTIMATE_URL = os.Getenv("QUICKBOOKS_DEBUG_DELETE_ESTIMATE_URL")
		
		if QUICKBOOKS_CLIENT_ID == "" || QUICKBOOKS_CLIENT_SECRET == "" || QUICKBOOKS_AUTH_CALLBACK_URL == "" || QUICKBOOKS_API_URL == "" || QUICKBOOKS_GET_CUSTOMER_URL == "" || QUICKBOOKS_GET_PRODUCT_URL == "" || QUICKBOOKS_CREATE_ESTIMATE_URL == "" || QUICKBOOKS_DELETE_ESTIMATE_URL == "" {
			log.Fatalf("Error initializing QuickBooks credentials: missing required environment variables")
		}
		log.Println("Initialized quickbooks credentials in debug...")
	})
}

// For Prod and Staging
func InitQuickBooksFromSecrets(ctx context.Context) {
	initQuickBooksOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		QUICKBOOKS_CLIENT_ID = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_CLIENT_ID")
		QUICKBOOKS_CLIENT_SECRET = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_CLIENT_SECRET")
		QUICKBOOKS_AUTH_CALLBACK_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_AUTH_CALLBACK_URL")
		QUICBOOK_AUTH_CALLBACK_REDIRECT_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_AUTH_CALLBACK_REDIRECT_URL")
		QUICKBOOKS_API_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_API_URL")
		QUICKBOOKS_WEBHOOK_VERIFY_TOKEN = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_WEBHOOK_VERIFY_TOKEN")
		QUICKBOOKS_GET_CUSTOMER_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_GET_CUSTOMER_URL")
		QUICKBOOKS_GET_PRODUCT_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_GET_PRODUCT_URL")
		QUICKBOOKS_CREATE_ESTIMATE_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_CREATE_ESTIMATE_URL")
		QUICKBOOKS_DELETE_ESTIMATE_URL = gcp.LoadSecretsHelper(projectID, "QUICKBOOKS_DELETE_ESTIMATE_URL")
		
		if QUICKBOOKS_CLIENT_ID == "" || QUICKBOOKS_CLIENT_SECRET == "" || QUICKBOOKS_AUTH_CALLBACK_URL == "" || QUICKBOOKS_API_URL == "" || QUICKBOOKS_WEBHOOK_VERIFY_TOKEN == "" || QUICKBOOKS_GET_CUSTOMER_URL == "" || QUICKBOOKS_GET_PRODUCT_URL == "" || QUICKBOOKS_CREATE_ESTIMATE_URL == "" || QUICKBOOKS_DELETE_ESTIMATE_URL == "" {
			log.Fatalf("Error initializing QuickBooks credentials: missing required environment variables")
		}
		log.Println("QuickBooks credentials initialized for PRODUCTION environment.")
	})
}
