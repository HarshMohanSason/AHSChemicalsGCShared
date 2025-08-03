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

// Company Details loaded at runtime.
var (
	COMPANYEMAIL           string            // Public-facing company email
	COMPANYPHONE           string            // Support or inquiry phone number
	COMPANYADDRESSLINE1    string            // Primary street address
	COMPANYADDRESSLINE2    string            // Secondary address line (e.g., Suite number)
	COMPANY24HOURPHONE     string            // 24-hour phone number
	EMAILINTERNALRECIPENTS map[string]string // Key-value map of internal email recipients
	LOGOPATH               string            // Path or URL to the company logo image
	initCompanyDetailsOnce sync.Once
)

// InitCompanyDetailsDebug initializes company configuration using environment
// variables for local development or debug environments.
// It logs fatal errors if any required variable is missing.
func InitCompanyDetailsDebug() {
	COMPANYPHONE = os.Getenv("COMPANYPHONE")
	COMPANYEMAIL = os.Getenv("COMPANYEMAIL")
	COMPANYADDRESSLINE1 = os.Getenv("COMPANYADDRESSLINE1")
	COMPANYADDRESSLINE2 = os.Getenv("COMPANYADDRESSLINE2")
	COMPANY24HOURPHONE = os.Getenv("COMPANY24HOURPHONE")
	rawEmailRecipients := os.Getenv("EMAILINTERNALRECIPIENTS")
	err := json.Unmarshal([]byte(rawEmailRecipients), &EMAILINTERNALRECIPENTS)
	if err != nil {
		log.Fatal(err)
	}

	LOGOPATH = os.Getenv("LOGOPATH")

	if COMPANYPHONE == "" || COMPANYEMAIL == "" || COMPANYADDRESSLINE1 == "" || COMPANYADDRESSLINE2 == "" || LOGOPATH == "" || COMPANY24HOURPHONE == "" {
		log.Fatal("Company details not initialized. Please check environment variables.")
	}

	log.Println("Company details initialized for DEBUG environment.")
}

// InitCompanyDetailsStaging initializes company configuration using Google Cloud Secret Manager
// It behaves similarly to InitCompanyDetailsProd and also logs fatal errors if any required variable is missing.
func InitCompanyDetailsStaging(ctx context.Context) {
	initCompanyDetailsOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		COMPANYPHONE = gcp.LoadSecretsHelper(projectID, "COMPANYPHONE")
		COMPANYEMAIL = gcp.LoadSecretsHelper(projectID, "COMPANYEMAIL")
		COMPANYADDRESSLINE1 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE1")
		COMPANYADDRESSLINE2 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE2")
		COMPANY24HOURPHONE = gcp.LoadSecretsHelper(projectID, "COMPANY24HOURPHONE")
		rawEmailRecipients := gcp.LoadSecretsHelper(projectID, "EMAILINTERNALRECIPIENTS")
		err = json.Unmarshal([]byte(rawEmailRecipients), &EMAILINTERNALRECIPENTS)
		if err != nil {
			log.Fatal(err)
		}

		LOGOPATH = gcp.LoadSecretsHelper(projectID, "LOGOPATH")

		if COMPANYPHONE == "" || COMPANYEMAIL == "" || COMPANYADDRESSLINE1 == "" || COMPANYADDRESSLINE2 == "" || LOGOPATH == "" || COMPANY24HOURPHONE == "" {
			log.Fatal("Company details not initialized. Please check secrets.")
		}

		log.Println("Company details initialized for PRODUCTION environment.")
	})
}

// InitCompanyDetailsProd initializes company configuration using
// Google Cloud Secret Manager in a production environment.
//
// It uses `gcp.LoadSecretsHelper` to fetch required secrets from GCP
// based on the current project ID. The values are loaded only once via sync.Once.
func InitCompanyDetailsProd(ctx context.Context) {
	initCompanyDetailsOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		COMPANYPHONE = gcp.LoadSecretsHelper(projectID, "COMPANYPHONE")
		COMPANYEMAIL = gcp.LoadSecretsHelper(projectID, "COMPANYEMAIL")
		COMPANYADDRESSLINE1 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE1")
		COMPANYADDRESSLINE2 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE2")
		COMPANY24HOURPHONE = gcp.LoadSecretsHelper(projectID, "COMPANY24HOURPHONE")
		rawEmailRecipients := gcp.LoadSecretsHelper(projectID, "EMAILINTERNALRECIPIENTS")
		err = json.Unmarshal([]byte(rawEmailRecipients), &EMAILINTERNALRECIPENTS)
		if err != nil {
			log.Fatal(err)
		}

		LOGOPATH = gcp.LoadSecretsHelper(projectID, "LOGOPATH")

		if COMPANYPHONE == "" || COMPANYEMAIL == "" || COMPANYADDRESSLINE1 == "" || COMPANYADDRESSLINE2 == "" || LOGOPATH == "" || COMPANY24HOURPHONE == "" {
			log.Fatal("Company details not initialized. Please check secrets.")
		}

		log.Println("Company details initialized for PRODUCTION environment.")
	})
}
