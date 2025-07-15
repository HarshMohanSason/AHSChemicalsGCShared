package company_details

import (
	"context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

// Basic Company information
const (
	COMPANYNAME       = "Azure Hospitality Supply"
	COMPANYWEBSITEURL = "https://azurehospitalitysupply.com"
)

// Company details
var (
	COMPANYEMAIL           string
	COMPANYPHONE           string
	COMPANYADDRESSLINE1    string
	COMPANYADDRESSLINE2    string
	EMAILINTERNALRECIPENTS []string
	LOGOPATH               string
)

func InitCompanyDetailsDebug() {
	COMPANYPHONE = os.Getenv("COMPANYPHONE")
	COMPANYEMAIL = os.Getenv("COMPANYEMAIL")
	COMPANYADDRESSLINE1 = os.Getenv("COMPANYADDRESSLINE1")
	COMPANYADDRESSLINE2 = os.Getenv("COMPANYADDRESSLINE2")
	EMAILINTERNALRECIPENTS = strings.Split(os.Getenv("EMAILINTERNALRECIPENTS"), ";")
	LOGOPATH = os.Getenv("LOGOPATH")
	if COMPANYPHONE == "" || COMPANYEMAIL == "" || COMPANYADDRESSLINE1 == "" || COMPANYADDRESSLINE2 == "" || LOGOPATH == "" {
		log.Fatal("Company details not initialized. Please check environment variables.")
	}
	log.Println("Company details initialized for DEBUG environment.")
}

func InitCompanyDetailsProd(ctx context.Context) {
	shared.InitCompanyDetails.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Error loading Google Cloud project ID: %v", err)
		}

		COMPANYPHONE = gcp.LoadSecretsHelper(projectID, "COMPANYPHONE")
		COMPANYEMAIL = gcp.LoadSecretsHelper(projectID, "COMPANYEMAIL")
		COMPANYADDRESSLINE1 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE1")
		COMPANYADDRESSLINE2 = gcp.LoadSecretsHelper(projectID, "COMPANYADDRESSLINE2")
		EMAILINTERNALRECIPENTS = strings.Split(gcp.LoadSecretsHelper(projectID, "EMAILINTERNALRECIPENTS"), ";")
		LOGOPATH = gcp.LoadSecretsHelper(projectID, "LOGOPATH")
		if COMPANYPHONE == "" || COMPANYEMAIL == "" || COMPANYADDRESSLINE1 == "" || COMPANYADDRESSLINE2 == "" || LOGOPATH == "" {
			log.Fatal("Company details not initialized. Please check environment variables.")
		}
		log.Println("Company details initialized for PRODUCTION environment.")
	})
}
