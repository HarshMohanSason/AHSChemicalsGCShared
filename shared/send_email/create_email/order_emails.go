package create_email

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/send_email"
)

// CreateOrderPlacedAdminEmail creates an email payload to notify internal admin recipients
// that a new order has been placed.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//   - attachments: slice of Attachment objects representing files to be attached to the email.(Purchase order pdf)
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for admin notification.
//
// Example:
//
//	order := &orders.Order{ /* filled order details */ }
//	attachments := CreateAttachments(base64Contents, mimeTypes, filenames)
//	emailData := CreateAdminOrderPlacedEmail(order, attachments)
func CreateOrderPlacedAdminEmail(order *models.Order, attachments []send_email.Attachment) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"customer_name":        order.Customer.Name,
			"order_number":         order.ID,
			"order_date":           order.CreatedAt.Format("January 2, 2006 at 3:04 PM"),
			"items":                createItemsDataForAdminEmail(order),
			"subtotal":             order.GetFormattedSubTotal(),
			"tax_rate":             order.GetFormattedTaxRate(),
			"tax_amount":           order.GetFormattedTaxAmount(),
			"total":                order.GetFormattedTotal(),
			"special_instructions": order.SpecialInstructions,
		},
		TemplateID:  send_email.ORDER_PLACED_ADMIN_TEMPLATE_ID,
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
func CreateOrderPlacedUserEmail(order *models.Order) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: map[string]string{order.Customer.Email: order.Customer.Name},
		Data: map[string]any{
			"order_number":         order.ID,
			"order_date":           order.CreatedAt.Format("January 2, 2006 at 3:04 PM"),
			"items":                createItemsDataForUserEmail(order),
			"special_instructions": order.SpecialInstructions,
		},
		TemplateID:  send_email.ORDER_PLACED_USER_TEMPLATE_ID,
		Attachments: []send_email.Attachment{},
	}
	return emailData
}

// CreateOrderStatusUpdatedAdminEmail creates an email payload to notify internal admin recipients
// that an order status status has been successfully updated.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for admin notification.
func CreateOrderStatusUpdatedAdminEmail(order *models.Order) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"order_number":     order.ID,
			"order_updated_at": order.UpdatedAt.Format("January 2, 2006 at 3:04 PM"),
			"order_status":     order.Status,
			"customer_name":    order.Customer.Name,
			"customer_email":   order.Customer.Email,
			"order_total":      order.GetFormattedTotal(),
			"items":            createItemsDataForUserEmail(order), //Only need the sku, desc and quantity for this admin email template
		},
		TemplateID:  send_email.ORDER_STATUS_UPDATED_ADMIN_TEMPLATE_ID,
		Attachments: []send_email.Attachment{},
	}
	return emailData
}

// CreateOrderStatusUpdatedUserEmail creates an email payload to notify the customer
// that their order status has been successfully updated.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for customer notification.
func CreateOrderStatusUpdatedUserEmail(order *models.Order) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: map[string]string{order.Customer.Email: order.Customer.Name},
		Data: map[string]any{
			"order_number":     order.ID,
			"order_status":     order.Status,
			"order_updated_at": order.UpdatedAt.Format("January 2, 2006 at 3:04 PM"),
		},
		TemplateID:  send_email.ORDER_STATUS_UPDATED_USER_TEMPLATE_ID,
		Attachments: []send_email.Attachment{},
	}
	return emailData
}

// CreateOrderItemsUpdatedAdminEmail creates an email payload to notify internal admin recipients
// that an order has been updated.
//
// Parameters:
//   - original: pointer to the Order object containing order and customer details.
//   - updated: pointer to the Order object containing order and customer details.
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for admin notification.
func CreateOrderItemsUpdatedAdminEmail(order *models.Order) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"order_number":     order.ID,
			"order_updated_at": order.UpdatedAt.Format("January 2, 2006 at 3:04 PM"),
			"order_status":     order.Status,
			"customer_name":    order.Customer.Name,
			"customer_email":   order.Customer.Email,
			"order_total":      order.GetFormattedTotal(),
			"items":            createItemsDataForAdminEmail(order),
		},
		TemplateID:  send_email.ORDER_ITEMS_UPDATED_ADMIN_TEMPLATE_ID,
		Attachments: []send_email.Attachment{},
	}
	return emailData
}

// CreateOrderItemsUpdatedUserEmail creates an email payload to notify the customer
// that their order has been updated.
//
// Parameters:
//   - original: pointer to the Order object containing order and customer details.
//   - updated: pointer to the Order object containing order and customer details.
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for customer notification.
func CreateOrderItemsUpdatedUserEmail(order *models.Order) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: map[string]string{order.Customer.Email: order.Customer.Name},
		Data: map[string]any{
			"customer_name":    order.Customer.Name,
			"order_number":     order.ID,
			"order_status":     order.Status,
			"items":            createItemsDataForUserEmail(order),
		},
		TemplateID:  send_email.ORDER_ITEMS_UPDATED_USER_TEMPLATE_ID,
		Attachments: []send_email.Attachment{},
	}
	return emailData
}

// CreateOrderDeliveredAdminEmail creates an email payload to notify internal admin recipients
// that an order has been delivered.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//   - attachments: slice of Attachment objects representing files to be attached to the email.(Shipping manifest and temporary invoice pdf)
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for admin notification.
func CreateOrderDeliveredAdminEmail(delivery *models.Delivery) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"order_number":  delivery.Order.ID,
			"customer_name": delivery.Order.Customer.Name,
			"items":         createItemsDataForAdminEmail(delivery.Order),
		},
		TemplateID:  send_email.ORDER_DELIVERED_ADMIN_TEMPLATE_ID,
	}
	return emailData
}

// CreateOrderDeliveredUserEmail creates an email payload to notify the customer
// that their order has been successfully delivered.
//
// Parameters:
//   - order: pointer to the Order object containing order and customer details.
//   - delivery: pointer to the DeliveryData object containing delivery details.
//
// Returns:
//   - *EmailMetaData: a pointer to the populated EmailMetaData object configured for customer notification.
func CreateOrderDeliveredUserEmail(delivery *models.Delivery) *send_email.EmailMetaData {
	emailData := &send_email.EmailMetaData{
		Recipients: map[string]string{delivery.Order.Customer.Email: delivery.Order.Customer.Name},
		Data: map[string]any{
			"order_number": delivery.Order.ID,
			"received_by":  delivery.ReceivedBy,
			"delivered_by": delivery.DeliveredBy,
			"delivered_at": delivery.DeliveredAt.Format("January 2, 2006 at 3:04 PM"),
		},
		TemplateID:  send_email.ORDER_DELIVERED_USER_TEMPLATE_ID,
	}
	return emailData
}
