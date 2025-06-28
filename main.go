package main

import (
	"log"
	"os"

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/joho/godotenv"
)

func main(){

	//Load the env file
	err := godotenv.Load("./keys/.env")
	if err != nil{
		log.Fatalf("Error occurred loading the env file %v", err)
	}

	//Load the firebase admin sdk's
	if(os.Getenv("ENV") == "DEBUG"){
	firebase_shared.InitFirebaseDebug(os.Getenv("FIREBASE_ADMIN_SDK_DEBUG"))
		log.Print("Initialized firebase debug admin sdk")
	}else{
		path := os.Getenv("FIREBASE_ADMIN_SDK_PROD")
		firebase_shared.InitFirebaseProd(&path)
		log.Print("Initialized firebase prod admin sdk")
	}
}