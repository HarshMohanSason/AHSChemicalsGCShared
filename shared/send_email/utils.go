package send_email

import (
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/orders"
)

// createItemsDataForAdminEmail returns a slice of maps containing item details 
// formatted for admin emails. Each item includes the SKU, a detailed description, 
// quantity, and total price (unit price * quantity). This representation is 
// suitable for internal communication where pricing information is required.
//
// Parameters:
//   - order: pointer to an Order struct containing order item details.
//
// Returns:
//   - []map[string]any: slice of mapped order items including price information.
func createItemsDataForAdminEmail(order *orders.Order) []map[string]any {
	orderItems := make([]map[string]any, 0)
	for _, item := range order.Items {
		mappedItem := make(map[string]any)
		mappedItem["sku"] = item.SKU
		mappedItem["description"] = fmt.Sprintf("%s-%s %.2f %s (Pack of %d)", item.Brand, item.Name, item.Size, item.SizeUnit, item.PackOf)
		mappedItem["quantity"] = item.Quantity
		mappedItem["price"] = fmt.Sprintf("$%.2f", item.UnitPrice*float64(item.Quantity))

		orderItems = append(orderItems, mappedItem)
	}
	return orderItems
}

// createItemsDataForUserEmail returns a slice of maps containing item details 
// formatted for user-facing emails. Each item includes the SKU, a detailed 
// description, and quantity. Price information is not included due to admin request
// or policy compliance.
//
// Parameters:
//   - order: pointer to an Order struct containing order item details.
//
// Returns:
//   - []map[string]any: slice of mapped order items excluding price information.
func createItemsDataForUserEmail(order *orders.Order) []map[string]any {
	orderItems := make([]map[string]any, 0)
	for _, item := range order.Items {
		mappedItem := make(map[string]any)
		mappedItem["sku"] = item.SKU
		mappedItem["description"] = fmt.Sprintf("%s-%s %.2f %s (Pack of %d)", item.Brand, item.Name, item.Size, item.SizeUnit, item.PackOf)
		mappedItem["quantity"] = item.Quantity

		orderItems = append(orderItems, mappedItem)
	}
	return orderItems
}

//CreateAttachments creates an array of Attachment structs to be sent along the SendGrid mail
//
//Parameters:
//  - base64Contents: []string
//  - mimeTypes: []string
//  - filenames: []string
//
//Returns:
//  - []Attachment
//
//Examples: 
//	 base64Contents := []string{
//	     "VGhpcyBpcyBmaWxlIGNvbnRlbnQ", // "File content 1"
//	 	 ...,
//	 }
//
//	 mimeTypes := []string{
//	     "text/plain",
//	     "application/pdf",
//	 }
//
//	 filenames := []string{
//	     "file1.txt",
//	     "file2.txt",
//	 }
//
//	 attachments := CreateAttachments(base64Contents, mimeTypes, filenames)
//
func CreateAttachments (base64Contents, mimeTypes, filenames []string) []Attachment{
	attachments := make([]Attachment, 0)
	for i := range base64Contents {
		attachment := Attachment{
			Base64Content: base64Contents[i],
			MimeType:      mimeTypes[i],
			FileName:      filenames[i],
		}
		attachments = append(attachments, attachment)
	}
	return attachments
}