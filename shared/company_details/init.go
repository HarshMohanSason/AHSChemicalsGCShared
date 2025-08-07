// Package company_details provides methods to initialize and access
// company-level details such as contact information, logo path, and internal recipients.
// Initialization varies by environment: Debug, Staging, or Production.
package company_details

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

var (
	COMPANYNAME            string            // Public-facing company name
	COMPANYURL             string            // Public-facing company URL
	COMPANYEMAIL           string            // Public-facing company email
	COMPANYPHONE           string            // Support or inquiry phone number
	COMPANYADDRESSLINE1    string            // Primary street address
	COMPANYADDRESSLINE2    string            // Secondary address line (e.g., Suite number)
	EMAILINTERNALRECIPENTS map[string]string // Key-value map of internal email recipients
	LOGOPATH               string            // Path or URL to the company logo image
	initCompanyDetailsOnce sync.Once
)

func InitCompanyDetailsDebug() {
	initCompanyDetailsOnce.Do(func() {
		COMPANYNAME = os.Getenv("COMPANYNAME")
		COMPANYURL = os.Getenv("COMPANYURL")
		COMPANYPHONE = os.Getenv("COMPANYPHONE")
		COMPANYEMAIL = os.Getenv("COMPANYEMAIL")
		COMPANYADDRESSLINE1 = os.Getenv("COMPANYADDRESSLINE1")
		COMPANYADDRESSLINE2 = os.Getenv("COMPANYADDRESSLINE2")
		rawEmailRecipients := os.Getenv("EMAILINTERNALRECIPIENTS")
		err := json.Unmarshal([]byte(rawEmailRecipients), &EMAILINTERNALRECIPENTS)
		if err != nil {
			log.Fatal(err)
		}
		LOGOPATH = os.Getenv("LOGOPATH")

		if COMPANYPHONE == "" || COMPANYEMAIL == "" || COMPANYADDRESSLINE1 == "" || COMPANYADDRESSLINE2 == "" || LOGOPATH == "" || COMPANYNAME == "" || COMPANYURL == "" {
			log.Fatal("Company details not initialized. Please check environment variables.")
		}

		log.Println("Company details initialized for DEBUG environment.")
	})
}

//For production and staging
func InitCompanyDetailsFromSecrets(ctx context.Context) {
	initCompanyDetailsOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		COMPANYNAME = gcp.LoadSecretsHelper(projectID, "COMPANYNAME")
		COMPANYURL = gcp.LoadSecretsHelper(projectID, "COMPANYURL")
		COMPANYPHONE = gcp.LoadSecretsHelper(projectID, "COMPANYPHONE")
		COMPANYEMAIL = gcp.LoadSecretsHelper(projectID, "COMPANYEMAIL")
		COMPANYADDRESSLINE1 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE1")
		COMPANYADDRESSLINE2 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE2")
		rawEmailRecipients := gcp.LoadSecretsHelper(projectID, "EMAILINTERNALRECIPIENTS")
		err = json.Unmarshal([]byte(rawEmailRecipients), &EMAILINTERNALRECIPENTS)
		if err != nil {
			log.Fatal(err)
		}
		LOGOPATH = gcp.LoadSecretsHelper(projectID, "LOGOPATH")

		if COMPANYPHONE == "" || COMPANYEMAIL == "" || COMPANYADDRESSLINE1 == "" || COMPANYADDRESSLINE2 == "" || LOGOPATH == "" || COMPANYNAME == "" || COMPANYURL == "" {
			log.Fatal("Company details not initialized. Please check secrets.")
		}

		log.Println("Company details initialized for PRODUCTION environment.")
	})
}
