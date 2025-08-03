package send_email

import (
	"context"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

var (
	SENDGRID_API_KEY string
	initSendGridOnce sync.Once
)

// Send grid dynamic email template IDs.
const (
	ACCOUNT_CREATED_USER_TEMPLATE_ID       = "d-7862af13a00340f4a252d405712ea368"
	ACCOUNT_DELETED_USER_TEMPLATE_ID       = "d-16455882508f49a69cc64af9df98bc79"
	CONTACT_US_ADMIN_TEMPLATE_ID           = "d-40d7a73d43044231ab2e3e20d6760db5"
	CONTACT_US_USER_TEMPLATE_ID            = "d-84ee771b07344758a5bd0fe38afd8ae8"
	ORDER_PLACED_ADMIN_TEMPLATE_ID         = "d-b3585038c1094054a8404423e6051573"
	ORDER_PLACED_USER_TEMPLATE_ID          = "d-5d7698237ee7491c96af123f964b4321"
	ORDER_STATUS_UPDATED_USER_TEMPLATE_ID  = "d-1deaf1f378e449ce919316337c0e1202"
	ORDER_STATUS_UPDATED_ADMIN_TEMPLATE_ID = "d-cd691ca20147468e82935bdc72303447"
	ORDER_ITEMS_UPDATED_USER_TEMPLATE_ID   = "d-f88446dbe2cb4b83b7452f89906163ee"
	ORDER_ITEMS_UPDATED_ADMIN_TEMPLATE_ID  = "d-2df971c9b243495faf98f119c7f25dd1"
	ORDER_DELIVERED_ADMIN_TEMPLATE_ID      = "d-10e63e3c62ac468da03efee2b15a5e96"
	ORDER_DELIVERED_USER_TEMPLATE_ID       = "d-a6344e4ac9254af28afde74ad2b81bc4"
)

// InitSendGridDebug initializes SendGrid credentials in debug mode
func InitSendGridDebug() {
	SENDGRID_API_KEY = os.Getenv("SENDGRID_API_KEY")
	if SENDGRID_API_KEY == "" {
		log.Fatal("SENDGRID_API_KEY is not set")
	}
	log.Println("Initialized SendGrid credentials in debug mode")
}

// InitSendGridStaging initializes SendGrid credentials in staging mode.
// Only used in the init functions of Google Cloud Functions for staging
//
// Parameters:
//   - context for the function
func InitSendGridStaging(ctx context.Context) {
	initSendGridOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}
		SENDGRID_API_KEY = gcp.LoadSecretsHelper(projectID, "SENDGRID_API_KEY")
		if SENDGRID_API_KEY == "" {
			log.Fatal("SENDGRID_API_KEY_STAGING is not set")
		}
		log.Println("Initialized SendGrid credentials in staging mode")
	})
}

// InitSendGridProd initializes SendGrid credentials in production mode.
// Only used in the init functions of google cloud functions
//
// Parameters:
//   - context for the function
func InitSendGridProd(ctx context.Context) {
	initSendGridOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}
		SENDGRID_API_KEY = gcp.LoadSecretsHelper(projectID, "SENDGRID_API_KEY")
		if SENDGRID_API_KEY == "" {
			log.Fatal("SENDGRID_API_KEY is not set")
		}
		log.Println("Initialized SendGrid credentials in production mode")
	})
}