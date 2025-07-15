package tests

import (
	"log"
	"os"
	"testing"

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M){

	err := godotenv.Load("../../../keys/.env")
	if err != nil{
		log.Print(err.Error())
	}
	log.Printf("Credentials: %s",os.Getenv("FIREBASE_ADMIN_SDK_DEBUG"))
	adminSDKFilePath := os.Getenv("FIREBASE_ADMIN_SDK_DEBUG")
	
	//Initialize the debug project sdk
	firebase_shared.InitFirebaseDebug(adminSDKFilePath)

	exitCode := m.Run()

	os.Exit(exitCode)
}