package tests

import (
	"github.com/joho/godotenv"
	"os"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
	"testing"
	)

func TestMain(m *testing.M){

	//Load the env file
	godotenv.Load("/Users/harshmohansason/Documents/AHSChemicalsGCShared/keys/.env")
	
	//Load the firebase debug project admin sdk
	debugPath := os.Getenv("FIREBASE_ADMIN_SDK_DEBUG")
	shared.InitFirebaseDebug(&debugPath)
	
	//Run 
	exitCode := m.Run()

	//Exit
	os.Exit(exitCode)
}