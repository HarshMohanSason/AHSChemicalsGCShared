package create_email

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/send_email"
)

// CreateAttachments creates an array of Attachment structs to be sent along the SendGrid mail
//
// Parameters:
//   - base64Contents: []string
//   - mimeTypes: []string
//   - filenames: []string
//
// Returns:
//   - []Attachment
//
// Examples:
//
//	base64Contents := []string{
//	    "VGhpcyBpcyBmaWxlIGNvbnRlbnQ", // "File content 1"
//		 ...,
//	}
//
//	mimeTypes := []string{
//	    "text/plain",
//	    "application/pdf",
//	}
//
//	filenames := []string{
//	    "file1.txt",
//	    "file2.txt",
//	}
//
//	attachments := CreateAttachments(base64Contents, mimeTypes, filenames)
//
// Note: filenames must have the extension with the name. In some applications such as outlook mobile, it will
// not render the attachment if the extension is not provided along with the filename
func CreateAttachments(base64Contents, mimeTypes, filenames []string) []send_email.Attachment {
	attachments := make([]send_email.Attachment, 0)
	for i := range base64Contents {
		attachment := send_email.Attachment{
			Base64Content: base64Contents[i],
			MimeType:      mimeTypes[i],
			FileName:      filenames[i],
		}
		attachments = append(attachments, attachment)
	}
	return attachments
}

func CreateSingleAttachment(base64Content, mimeType, filename string) send_email.Attachment {
	return send_email.Attachment{
		Base64Content: base64Content,
		MimeType:      mimeType,
		FileName:      filename,
	}
}

// createItemsDataForAdminEmail returns a slice of maps containing item details
// formatted for admin emails. Each item includes the SKU, a detailed description,
// quantity, and total price (unit price * quantity). This data is used to render the 
// items list used to display in the dynamic send grid emails ( {{#each items}}...)
//
// Parameters:
//   - order: pointer to an Order struct containing order item details.
//
// Returns:
//   - []map[string]any: slice of mapped order items including price information.
func createItemsDataForAdminEmail(order *models.Order) []map[string]any {
	orderItems := make([]map[string]any, 0)
	for _, item := range order.Items {
		mappedItem := make(map[string]any)
		mappedItem["sku"] = item.SKU
		mappedItem["description"] = item.GetFormattedDescription()
		mappedItem["quantity"] = item.Quantity
		mappedItem["price"] = item.GetFormattedTotalPrice()

		orderItems = append(orderItems, mappedItem)
	}
	return orderItems
}

func createItemsDetailedItemsDataForAdminEmail(order *models.Order) []map[string]any{
	orderItems := make([]map[string]any, 0)
	for _, item := range order.Items {
		mappedItem := make(map[string]any)
		mappedItem["sku"] = item.SKU
		mappedItem["description"] = item.GetFormattedDescription()
		mappedItem["quantity"] = item.Quantity
		mappedItem["unit_price"] = item.GetFormattedUnitPrice()
		mappedItem["total_price"] = item.GetFormattedTotalPrice()

		orderItems = append(orderItems, mappedItem)
	}
	return orderItems
}

// createItemsDataForUserEmail returns a slice of maps containing item details
// formatted for user-facing emails. Each item includes the SKU, a detailed
// description, and quantity. Price information is not included due to owner request.
// This data is used to render the items list used to display in the dynamic send grid emails ( {{#each items}}...)
//
// Parameters:
//   - order: pointer to an Order struct containing order item details.
//
// Returns:
//   - []map[string]any: slice of mapped order items excluding price information.
func createItemsDataForUserEmail(order *models.Order) []map[string]any {
	orderItems := make([]map[string]any, 0)
	for _, item := range order.Items {
		mappedItem := make(map[string]any)
		mappedItem["sku"] = item.SKU
		mappedItem["description"] = item.GetFormattedDescription()
		mappedItem["quantity"] = item.Quantity

		orderItems = append(orderItems, mappedItem)
	}
	return orderItems
}