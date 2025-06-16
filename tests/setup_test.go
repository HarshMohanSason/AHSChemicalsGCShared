package tests

import (
	"log"
	"os"
	"testing"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M){

	//Load the env file
	err := godotenv.Load("../keys/.env")
	if err != nil{
		log.Fatalf("Error occurred loading the env file %v", err)
	}

	//Load the admin sdk json from env
	debugPath := os.Getenv("FIREBASE_ADMIN_SDK_DEBUG")
	
	//Initialize the debug project
	shared.InitFirebaseDebug(debugPath)
	
	//Run 
	exitCode := m.Run()

	//Exit
	os.Exit(exitCode)
}