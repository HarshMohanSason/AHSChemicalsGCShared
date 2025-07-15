package send_email

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

var (
	SENDGRID_API_KEY   string 
)

func InitSendGridDebug() {
	SENDGRID_API_KEY = os.Getenv("SENDGRID_API_KEY")
	if (SENDGRID_API_KEY == "") {
		log.Fatal("SENDGRID_API_KEY is not set")
	}
	log.Println("Initialized SendGrid credentials in debug mode")
}

func InitSendGridProd(ctx context.Context) {
	shared.InitSendGridOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}
		SENDGRID_API_KEY = gcp.LoadSecretsHelper(projectID, "SENDGRID_API_KEY")
		if (SENDGRID_API_KEY == "") {
			log.Fatal("SENDGRID_API_KEY is not set")
		}
		log.Println("Initialized SendGrid credentials in production mode")
	})
}