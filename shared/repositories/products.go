package repositories

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
)

//FetchAllProductsByIDs fetches all products from firestore from collection ('products')
//Returns map of product id and the product object in order to search in O(1)
//
//Params:
//  - ctx: context
//  - productIDs: []string, product ids
//
//Returns:
//  - map[string]models.Product. Key is product id
//  - error
func FetchAllProductsByIDs(ctx context.Context, productIDs []string) (map[string]models.Product, error) {
	docRefs := make([]*firestore.DocumentRef, len(productIDs))
	for i, productID := range productIDs {
		docRefs[i] = firebase_shared.FirestoreClient.Collection("products").Doc(productID)
	}
	docSnapshots, err := firebase_shared.FirestoreClient.GetAll(ctx, docRefs)
	if err != nil{
		return nil, err
	}

	productMap := make(map[string]models.Product)
	for i, docSnapshot := range docSnapshots {
		var product models.Product
		if err := docSnapshot.DataTo(&product); err != nil {
			return nil, fmt.Errorf("error decoding product %s: %v", productIDs[i], err)
		}
		productMap[product.ID] = product
	}
	return productMap, nil
}

// FetchAllProductsFromFirestore fetches all products from firestore from collection ('products')
func FetchAllProductsFromFirestore(ctx context.Context) ([]models.Product, error) {

	products, err := firebase_shared.FirestoreClient.Collection("products").Where("isActive", "==", true).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	productList := make([]models.Product, len(products))
	for i, product := range products {
		var item models.Product
		if err := product.DataTo(&item); err != nil {
			return nil, fmt.Errorf("error decoding product %s: %v", product.Ref.ID, err)
		}
		productList[i] = item
	}
	return productList, nil
}

// SyncQuickbookProductRespToFirestore syncs quickbook product response to firestore
// collection ('products')
//
// Params:
//   - qbItemsResponse: *qbmodels.QBItemsResponse, a mapped response from quickbooks
//   - ctx: context
//
// Returns:
//   - error: error
func SyncQuickbookProductRespToFirestore(qbItemsResponse *qbmodels.QBItemsResponse, ctx context.Context) error {
	if qbItemsResponse == nil || qbItemsResponse.QueryResponse.Item == nil {
		return nil
	}

	bulkWriter := firebase_shared.FirestoreClient.BulkWriter(ctx)

	for _, item := range qbItemsResponse.QueryResponse.Item {
		docRef := firebase_shared.FirestoreClient.Collection("products").Doc(item.ID)
		_, err := bulkWriter.Set(docRef, item.MapToProduct().ToMap(), firestore.MergeAll)
		if err != nil {
			return err
		}
	}
	bulkWriter.Flush()
	return nil
}

// FetchProductFromFirestore fetches a single product from firestore. 
//
// Params:
//   - ctx: context
//   - productID: string, product id
//
// Returns:
//   - *models.Product, error
func FetchProductFromFirestore(ctx context.Context, productID string) (*models.Product, error) {
	doc, err := firebase_shared.FirestoreClient.Collection("products").Doc(productID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var product models.Product
	err = doc.DataTo(&product)
	if err != nil{
		return nil, err
	}
	return &product, nil
}

//UpdateProductInFirestore updates current product document in firestore. 
//
// Params: 
//  - ctx context of the async function
//  - productID is the id of the product whose doc is being updated
//  - details is of type any, the fields that need to be updated
// 
// Returns: 
//  - error: if any 
// Note: The details is of type any because any field can be updated when called. So its much
// better to just pass the type any. In Addition to this, before calling this, make sure to always
// and always double check if the key matches with the `firestore` key in the struct otherwise this
// will create a new key with that value in document. 
func UpdateProductInFirestore(ctx context.Context, productID string, details any) error {
	_, err := firebase_shared.FirestoreClient.Collection("products").Doc(productID).Set(ctx, details, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}


// CreateProductInFirestore creates new product document in firestore.
//
// Params:
//   - ctx: context
//   - customer: *models.Customer
//
// Returns:
//   - error: error
func CreateProductInFirestore(ctx context.Context, product *models.Product) error {
	_, err := firebase_shared.FirestoreClient.Collection("products").Doc(product.ID).Set(ctx, product.ToMap())
	return err
}