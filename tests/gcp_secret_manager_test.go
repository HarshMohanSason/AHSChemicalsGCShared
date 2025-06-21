package tests

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/compute/metadata"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
)

func TestGetSecretFromGCP(t *testing.T){
	
	projectID, err := metadata.ProjectIDWithContext(context.Background())
	if err != nil{
		t.Fatalf("Error getting the project id %v", err)
	}

	secretPath := fmt.Sprintf("projects/%s/secrets/%s/versions/latest",projectID,"TWILIO_RECIPENTS_PHONE")

	_, err = shared.GetSecretFromGCP(secretPath)
	if err != nil{
		t.Fatalf("Error occurred while getting the secret %v", err)
	}
}