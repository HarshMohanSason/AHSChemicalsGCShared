package main

import (
	"context"
	"log"
	"os"
	"time"

	//firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/customers"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/orders"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/purchase_order"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/products"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("./keys/.env")
	if err != nil {
		log.Fatalf("Error occurred loading the env file %v", err)
	}

	if os.Getenv("ENV") == "DEBUG" {
		company_details.InitCompanyDetailsDebug()
		//firebase_shared.InitFirebaseDebug(os.Getenv("FIREBASE_ADMIN_SDK_DEBUG"))
	} else {
		company_details.InitCompanyDetailsProd(context.Background())
		//path := os.Getenv("FIREBASE_ADMIN_SDK_PROD")
		//firebase_shared.InitFirebaseProd(&path)
	}

	//Test order
	order := &orders.Order{
		ID: "112312ddsfad",
		Customer: customers.Customer{
			ID:          "112312ddsfad",
			DisplayName: "Harsh",
			PrimaryEmailAddr: &customers.Email{
				Address: "QwKu9@example.com",
			},
			PrimaryPhone: &customers.Phone{
				FreeFormNumber: "1231222222",
			},
			BillAddr: &customers.Address{
				Line1:                  "2040 N Preisker lane",
				City:                   "Santa Maria",
				CountrySubDivisionCode: "CA",
				PostalCode:             "93726",
			},
		},
		SpecialInstructions: "Testing this first product. A quick brown fox jumps over the lazy dog. Test speical instrucitons",
		Status:              "PENDING", // basic default
		Items: []products.Item{
			{
				Name:        "Test product1 ASASDJ ASD ASJ DASJ DJSA DASDJASJDAS. DJS ",
				ID:          "112312ddsfad",
				Quantity:    10,
				UnitPrice:   5.0,
				Brand:       "MicroTECH",
				SKU:         "asef",
				SizeUnit:    "GAL",
				Size:        12.0,
				Description: "Test product description",
			},
			{
				Name:        "Test product2",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			{
				Name:        "Test product3",
				ID:          "112312ddsfad",
				Quantity:    3,
				UnitPrice:   2.0,
				Brand:       "MicroTECH",
				SKU:         "PB0012612AB",
				SizeUnit:    "GAL",
				Size:        10.0,
				Description: "Test product description2",
			},
			


		},
		TaxRate:   0.1253,
		SubTotal:  1232.123,
		TaxAmount: 12123.123,
		Total:     12312.123,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	purchase_order.CreatePurchaseOrderPDF(order)
}
