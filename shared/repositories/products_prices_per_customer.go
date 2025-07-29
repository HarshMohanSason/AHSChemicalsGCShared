package repositories

import (
	"context"
	"errors"

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

func SaveProductsPricesPerCustomerToFirestore(ctx context.Context) error {

	customers, err := FetchAllCustomersFromFirestore(ctx)
	if err != nil {
		return err
	}
	products, err := FetchAllProductsFromFirestore(ctx)
	if err != nil {
		return err
	}
	totalExpectedDocs := len(customers) * len(products)

	productPricesPerCustomer, err := firebase_shared.FirestoreClient.Collection("products_prices_per_customer").Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	if len(productPricesPerCustomer) == totalExpectedDocs {
		return errors.New("No new customer or product found to sync product prices for customers")
	}

	counter := 0
	bulkWriter := firebase_shared.FirestoreClient.BulkWriter(ctx)
	
	for _, product := range products {
		for _, customer := range customers {
			productPricesDocs, err := firebase_shared.FirestoreClient.Collection("products_prices_per_customer").
				Where("productId", "==", product.ID).
				Where("customerId", "==", customer.ID).
				Documents(ctx).GetAll()
			if err != nil {
				return err
			}
			//If no docs found then create one for the customer
			if len(productPricesDocs) == 0 {
				productToAdd := models.CreateProductPricePerCustomer(&product, customer.ID).ToMap()
				_, err := bulkWriter.Create(firebase_shared.FirestoreClient.Collection("products_prices_per_customer").NewDoc(), productToAdd)
				if err != nil {
					return err
				}
			}
			counter++; 
			if counter % 500 == 0 {	
				bulkWriter.Flush()
				bulkWriter = firebase_shared.FirestoreClient.BulkWriter(ctx)
			}
		}
	}
	// Flush remaining docs
	if counter % 500 != 0 {
		bulkWriter.Flush()
	}
	return nil
}
