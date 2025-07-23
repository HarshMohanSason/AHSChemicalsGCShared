package send_email

import (
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/orders"
)

// CreateOrderPlacedAdminEmail creates an email payload to notify internal admin recipients
// that a new order has been placed.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//   - attachments: slice of Attachment objects representing files to be attached to the email.
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for admin notification.
//
// Example:
//
//   order := &orders.Order{ /* filled order details */ }
//   attachments := CreateAttachments(base64Contents, mimeTypes, filenames)
//   emailData := CreateAdminOrderPlacedEmail(order, attachments)
//
func CreateOrderPlacedAdminEmail(order *orders.Order, attachments []Attachment) *EmailMetaData {
	emailData := &EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"customer_name": order.Customer.DisplayName,
			"order_number": order.ID,
			"order_date": order.CreatedAt.Format("January 2, 2006 at 3:04 PM MST"),
			"items": createItemsDataForAdminEmail(order),
			"subtotal": fmt.Sprintf("$%.2f", order.SubTotal),
			"tax_rate": fmt.Sprintf("%.2f%%", order.TaxRate* 100),
			"tax_amount": fmt.Sprintf("$%.2f", order.TaxAmount),
			"total": fmt.Sprintf("$%.2f", order.Total),
			"special_instructions": order.SpecialInstructions,
		},
		TemplateID: ORDER_PLACED_ADMIN_TEMPLATE_ID,
		Attachments: attachments,
	}
	return emailData
}

// CreateOrderPlacedUserEmail creates an email payload to notify the customer
// that their order has been successfully placed.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for customer notification.
//
// Example:
//
//   order := &orders.Order{ /* filled order details */ }
//   emailData := CreateUserOrderPlacedEmail(order)
//
func CreateOrderPlacedUserEmail(order *orders.Order) *EmailMetaData {
	emailData := &EmailMetaData{
		Recipients: map[string]string{order.Customer.PrimaryEmailAddr.Address: order.Customer.DisplayName},
		Data: map[string]any{
			"order_number": order.ID,
			"order_date": order.CreatedAt.Format("January 2, 2006 at 3:04 PM MST"),
			"items": createItemsDataForUserEmail(order),
			"special_instructions": order.SpecialInstructions,
		},
		TemplateID: ORDER_PLACED_USER_TEMPLATE_ID,
		Attachments: []Attachment{},
	}
	return emailData
}

// CreateOrderUpdatedUserEmail creates an email payload to notify the customer
// that their order has been successfully updated.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//
// Returns: 
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for customer notification.
//
func CreateOrderUpdatedUserEmail(order *orders.Order) *EmailMetaData {
	emailData := &EmailMetaData{
		Recipients: map[string]string{order.Customer.PrimaryEmailAddr.Address: order.Customer.DisplayName},
		Data: map[string]any{
			"order_number": order.ID,
			"order_updated_at": order.UpdatedAt.Format("January 2, 2006 at 3:04 PM MST"),
			"order_status":order.Status,
			"items": createItemsDataForUserEmail(order),
		},
		TemplateID: ORDER_UPDATED_USER_TEMPLATE_ID,
		Attachments: []Attachment{},
	}
	return emailData
}