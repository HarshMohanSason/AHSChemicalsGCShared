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
	
	//Load the admin sdk json from env
	debugPath := os.Getenv("FIREBASE_ADMIN_SDK_DEBUG")
	
	//Initialize the debug project
	shared.InitFirebaseDebug(debugPath)
	
	//Run 
	exitCode := m.Run()

	//Exit
	os.Exit(exitCode)
}