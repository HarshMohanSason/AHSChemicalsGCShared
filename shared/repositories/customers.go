package repositories

import (
	"context"
	"fmt"

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

// FetchCustomerFromFirestore fetches customer object from firestore collection ('customers')
// based on customerID
//
// Params:
// 	- customerID: string
// 	- ctx: context
//
// Returns:
//  - *firestore_models.Customer
//  - error: error
//
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