package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	
	err := godotenv.Load("../../../keys/.env.development")
	if err != nil{
		log.Fatalf("Error loading the .env file: %v", err)
	}
	gcp.InitPubSubDebug(context.Background())
	code := m.Run()

	os.Exit(code)
}