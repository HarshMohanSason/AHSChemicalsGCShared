package main

import (
	"log"
	"os"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
	"github.com/joho/godotenv"
)

func main(){

	//Load the env file
	err := godotenv.Load("keys/.env")
	if err != nil{
		log.Fatalf("Error occurred loading the env file %v", err)
	}
	//Load the admin sdk
	if(os.Getenv("ENV") == "DEBUG"){
		shared.InitFirebaseDebug(os.Getenv("FIREBASE_ADMIN_SDK_DEBUG"))
		log.Print("Initialized firebase debug admin sdk")
	}else{
		path := os.Getenv("FIREBASE_ADMIN_SDK_PROD")
		shared.InitFirebaseProd(&path)
		log.Print("Initialized firebase prod admin sdk")
	}

	//Call any functions to test here locally
}