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

//GetProductPricesFromCustomerID returns a map of product id to price. 
//Returning a map is done to avoid multiple calls to firestore and to search in O(1)
//
//Params:
//	- customerID: ID of the customer
//	- ctx: context
//
//Returns:
//	- map of product id to price
//	- error
func GetProductPricesFromCustomerID(customerID string, ctx context.Context) (map[string]float64, error){

	docs, err := firebase_shared.FirestoreClient.Collection("products_prices_per_customer").Where("customerId", "==", customerID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	pricesMap := make(map[string]float64) //Map of product id to price
	for _, doc := range docs {
		var product models.ProductPricePerCustomer
		err := doc.DataTo(&product)
		if err != nil {
			return nil, err
		}
		pricesMap[product.ProductID] = product.Price
	}
	return pricesMap, nil
}
