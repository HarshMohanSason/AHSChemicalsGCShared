package repositories

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
)

// FetchCustomerFromFirestore fetches customer object from firestore collection ('customers')
// based on customerID
//
// Params:
//   - customerID: string
//   - ctx: context
//
// Returns:
//   - *firestore_models.Customer
//   - error: error
func FetchCustomerFromFirestore(id string, ctx context.Context) (*models.Customer, error) {

	docSnapshot, err := firebase_shared.FirestoreClient.Collection("customers").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var customer models.Customer
	err = docSnapshot.DataTo(&customer)
	if err != nil {
		return nil, fmt.Errorf("Error getting customer: %v", err)
	}
	return &customer, nil
}

// Only fetches active customers from firestore collection ('customers')
func FetchAllCustomersFromFirestore(ctx context.Context) ([]models.Customer, error) {

	customers, err := firebase_shared.FirestoreClient.Collection("customers").Where("isActive", "==", true).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	customersList := make([]models.Customer, len(customers))
	for i, customer := range customers {
		var customerObj models.Customer
		err = customer.DataTo(&customerObj)
		if err != nil {
			return nil, fmt.Errorf("Error getting customer: %v", err)
		}
		customersList[i] = customerObj
	}
	return customersList, nil
}

// SyncQuickbookCustomerRespToFirestore syncs quickbook customer response to firestore
// collection ('customers')
//
// Params:
//   - qbItemsResponse: *qbmodels.QBCustomersResponse, a mapped response from quickbooks
//   - ctx: context
//
// Returns:
//   - error: error
func SyncQuickbookCustomerRespToFirestore(qbItemsResponse *qbmodels.QBCustomersResponse, ctx context.Context) error {
	if qbItemsResponse == nil || qbItemsResponse.QueryResponse.Customer == nil {
		return nil
	}

	bulkWriter := firebase_shared.FirestoreClient.BulkWriter(ctx)

	for _, customer := range qbItemsResponse.QueryResponse.Customer {
		docRef := firebase_shared.FirestoreClient.Collection("customers").Doc(customer.ID)
		_, err := bulkWriter.Set(docRef, customer.MapToCustomer().ToMap(), firestore.MergeAll)
		if err != nil {
			return err
		}
	}
	bulkWriter.Flush()
	return nil
}

//UpdateCustomerInFirestore updates current customer document in firestore. 
//
// Params: 
//  - ctx context of the async function
//  - customerID is the id of the customer whose doc is being updated
//  - details is of type any, the fields that need to be updated
// 
// Returns: 
//  - error: if any 
// Note: The details is of type any because any field can be updated when called. So its much
// better to just pass the type any. In Addition to this, before calling this, make sure to always
// and always double check if the key matches with the `firestore` key in the struct otherwise this
// will create a new key with that value in document. 
func UpdateCustomerInFirestore(ctx context.Context, customerID string, details any) error {
	_, err := firebase_shared.FirestoreClient.Collection("customers").Doc(customerID).Set(ctx, details, firestore.MergeAll)

	if err != nil {
		return err
	}
	return nil
}

// CreateCustomerInFirestore creates new customer document in firestore.
//
// Params:
//   - ctx: context
//   - customer: *models.Customer
//
// Returns:
//   - error: error
//
func CreateCustomerInFirestore(ctx context.Context, customer *models.Customer) error {
	_, err := firebase_shared.FirestoreClient.Collection("customers").Doc(customer.ID).Set(ctx, customer.ToMap())
	return err
}