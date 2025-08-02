package repositories

import (
	"context"
	"encoding/base64"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/constants"
	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

// CreateOrderInFirestore creates a new order document in the Firestore "orders" collection.
//
// It takes a pointer to an Order object and a context. The function uploads the order data
// to Firestore using FirestoreClient. An auto-generated document ID is returned by Firestore
// and is then assigned back to the order's ID field.
//
// Parameters:
//   - order: Pointer to the Order struct containing the order details.
//   - ctx: Context used for the Firestore operation (should include timeout/deadline if needed).
//
// Returns:
//   - error: Returns an error if the Firestore insertion fails; otherwise, returns nil.
func CreateOrderInFirestore(order *models.Order, ctx context.Context) error {

	docRef, _, err := firebase_shared.FirestoreClient.Collection("orders").Add(ctx, order.ToMap())
	if err != nil {
		return err
	}
	order.SetID(docRef.ID)
	return nil
}

// CanPlaceOrder checks if there is already an order pending for the user if the status
// of any of the existing orders is still 'PENDING'. Firestore custom claims should be checked generally before calling this
// function to make sure it only checks for user since admins can place as many orders 
// as they want
//
// Parameters: 
// - orderUID - user id of of the order
// - ctx - context for the async function
// 
// Returns: 
// - error if any
func CanPlaceOrder(orderUID string, ctx context.Context) error {
	docs, err := firebase_shared.FirestoreClient.Collection("orders").Where("status", "==", constants.OrderStatusPending).Where("uid", "==", orderUID).Documents(ctx).GetAll()
	if err != nil{
		return err
	}
	if len(docs) > 0 {
		return fmt.Errorf("There is already an order pending for your account. Please contact the admin to change the status of the order to either approved or rejected in order to place a new order")
	}
	return nil
}

// UploadOrderFileToStorage uploads a base64-encoded PDF file to Cloud Storage
// under the path "orders/{orderID}/{fileName}".
//
// Parameters:
//   - orderID: the Firestore order document ID
//   - base64Str: base64 string of the file contents
//   - fileName: the name of the file to be stored (e.g. "invoice.pdf")
//   - ctx: the context for the upload operation
//
// Returns:
//   - error: if the upload fails at any step
func UploadOrderFileToStorage(orderID string, base64Str string, fileName string, ctx context.Context) error {
	bucket, err := firebase_shared.StorageClient.Bucket(firebase_shared.StorageBucket)
	if err != nil {
		return err
	}
	object := bucket.Object(fmt.Sprintf("orders/%s/%s", orderID, fileName))
	writer := object.NewWriter(ctx)

	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("failed to decode base64 string: %w", err)
	}

	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("failed to write to Cloud Storage: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

// UpdateOrderInFirestore updates specific fields of an order in Firestore.
// 
// Parameters:
// - orderID: Firestore document ID of the order to be updated.
// - details: A struct or map containing only the fields to be updated.
// - ctx: Context for Firestore operations.
//
// Note:
// - This function uses `firestore.MergeAll` to update only the specified fields in the existing document.
// - Pass only the necessary fields in `details`. Passing the entire order object may cause errors if it includes
//   fields that are not mapped correctly or are excluded with the `firestore:"-"` tag.
func UpdateOrderInFirestore(orderID string, details any, ctx context.Context) error {
	_, err := firebase_shared.FirestoreClient.Collection("orders").Doc(orderID).Set(ctx, details, firestore.MergeAll)

	if err != nil {
		return err
	}
	return nil
}

// FetchOrderFromFirestore fetches a basic order document from Firestore using the provided order ID.
//
// Parameters:
// - orderID: Firestore document ID of the order to fetch.
// - ctx: Context for Firestore operations.
//
// Returns:
// - A pointer to the `Order` struct if the document is found and successfully parsed.
// - An error if the document doesn't exist or the parsing fails.
func FetchOrderFromFirestore(orderID string, ctx context.Context) (*models.Order, error) {

	docSnapshot, err := firebase_shared.FirestoreClient.Collection("orders").Doc(orderID).Get(ctx)
	if err != nil {
		return nil, err
	}
	if !docSnapshot.Exists() {
		return nil, fmt.Errorf("order with ID %s not found", orderID)
	}
	var order models.Order
	if err := docSnapshot.DataTo(&order); err != nil {
		return nil, err
	}
	order.SetID(docSnapshot.Ref.ID)
	return &order, nil
}

// FetchDetailedOrderFromFirestore fetches an order document and enriches it with detailed data.
//
// Description:
// - Retrieves the order document from Firestore.
// - Populates the embedded `Customer` object using the `customerId` field in the order.
// - Converts the `items` field (a list of product IDs) into full product objects.
//
// Parameters:
// - orderID: Firestore document ID of the order to fetch.
// - ctx: Context for Firestore operations.
//
// Returns:
// - A pointer to the enriched `Order` struct containing full customer and product information.
// - An error if any part of the fetch or conversion fails.
func FetchDetailedOrderFromFirestore(orderID string, ctx context.Context) (*models.Order, error) {
	docSnapshot, err := firebase_shared.FirestoreClient.Collection("orders").Doc(orderID).Get(ctx)
	if err != nil {
		return nil, err
	}
	if !docSnapshot.Exists() {
		return nil, fmt.Errorf("order with ID %s not found", orderID)
	}
	//Assign the order object
	var order models.Order
	if err := docSnapshot.DataTo(&order); err != nil { 
		return nil, err
	}
	order.SetID(docSnapshot.Ref.ID)
	
	customerID := docSnapshot.Data()["customerId"].(string)
	customer, err := FetchCustomerFromFirestore(customerID, ctx)
	if err != nil {
		return nil, err
	}
	order.Customer = *customer

	//Get the product map
	productMap, err := FetchAllProductsByIDs(ctx, order.ToProductIDs())
	if err != nil{
		return nil, err
	}
	order.ToCompleteOrderItemsFromMinimal(productMap)

	return &order, nil
}