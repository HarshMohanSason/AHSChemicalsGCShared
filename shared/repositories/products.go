package repositories

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
)

// FetchProductByIDs returns a list of firestore products from a list of ids
//
// Parameters:
//   - ids []string: List of ids of the products to fetch
//   - ctx context.Context: context for the request
//
// Returns:
//   - []firestore_models.Product: List of products
//   - error:
func FetchProductByIDs(ids []string, ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	// Get all document references for customers collection
	collection := firebase_shared.FirestoreClient.Collection("customers")
	refs := make([]*firestore.DocumentRef, 0, len(ids))
	for _, id := range ids {
		refs = append(refs, collection.Doc(id))
	}

	docSnapshots, err := firebase_shared.FirestoreClient.GetAll(ctx, refs)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %v", err)
	}

	for _, doc := range docSnapshots {
		if !doc.Exists() {
			log.Printf("No document exists for the user when calling FetchProductByIDs: %s", doc.Ref.ID)
			continue // Moving on to next iteration if the document doesn't exist
		}
		var item models.Product
		if err := doc.DataTo(&item); err != nil {
			return nil, fmt.Errorf("error decoding customer %s: %v", doc.Ref.ID, err)
		}
		products = append(products, item)
	}

	return products, nil
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
