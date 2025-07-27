package firestore_repositories

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/services"
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
	log.Printf("Creating an order in firestore for the user: %s", order.Uid)

	docRef, _, err := firebase_shared.FirestoreClient.Collection("orders").Add(ctx, order)
	if err != nil {
		return err
	}

	order.ID = docRef.ID // Assign the generated Firestore ID back to the order object
	return nil
}

// GetCorrectCustomerPriceForProducts retrieves the custom prices for a specific customer
// and maps them to the corresponding product IDs.
//
// Parameters:
//   - order: the order containing the customer information
//   - ctx: the context for the Firestore operation
//
// Returns:
//   - map[string]float64: a map of product IDs to customer-specific prices
//   - error: if the Firestore query fails or data is malformed
func GetCorrectCustomerPriceForProducts(order *models.Order, ctx context.Context) (map[string]float64, error) {
	log.Printf("Fetching customer-specific prices for customer ID: %v", order.Customer.ID)

	snapshots, err := firebase_shared.FirestoreClient.Collection("product_prices").
		Where("customer_id", "==", order.Customer.ID).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customer prices: %w", err)
	}

	priceMap := make(map[string]float64)
	for _, doc := range snapshots {
		data := doc.Data()
		productID, _ := data["product_id"].(string)
		price, _ := data["price"].(float64)
		priceMap[productID] = price
	}

	return priceMap, nil
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
	log.Printf("Uploading file %s for order %s", fileName, orderID)

	bucket, err := firebase_shared.StorageClient.Bucket(firebase_shared.StorageBucket)
	if err != nil {
		log.Printf("Failed to get Cloud Storage bucket: %v", err)
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
	log.Printf("Updating an order in firestore for the admin")

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

	log.Printf("Fetching an order from firestore")

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
	log.Printf("Fetching an order from firestore")

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

	// Assign the customer object
	customerID := docSnapshot.Data()["customerId"].(string)
	customer, err := FetchCustomerFromFirestore(customerID, ctx)
	if err != nil {
		return nil, err
	}
	order.Customer = *customer

	// Assign the order items
	rawItems, ok := docSnapshot.Data()["items"].([]any)
	if !ok {
		return nil, fmt.Errorf("items field not in expected format")
	}
	orderIDs := services.ExtractOrderIDsFromFirestoreOrders(rawItems)
	orderItems, err := FetchProductByIDs(orderIDs, ctx)
	if err != nil {
		return nil, err
	}
	order.Items = orderItems

	return &order, nil
}