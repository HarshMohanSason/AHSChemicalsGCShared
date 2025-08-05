package create_email

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/send_email"
)

func CreateQuickBooksSessionExpiredEmail(reconnectURL string) *send_email.EmailMetaData{
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"reconnect_url": reconnectURL,
		},
		TemplateID:  send_email.QUICKBOOKS_SESSION_EXPIRED_TEMPLATE_ID,
		Attachments: []send_email.Attachment{},
	}
	return emailData
}

func CreateQuickBooksInvoiceAdminEmail(order *models.Order, invoice *qbmodels.Invoice) *send_email.EmailMetaData{
	emailData := &send_email.EmailMetaData{
		Recipients: company_details.EMAILINTERNALRECIPENTS,
		Data: map[string]any{
			"qb_invoice_id": invoice.ID,
			"order_number": order.ID,
			"customer_name": order.Customer.Name,
			"invoice_date": invoice.TxnDate,
			"items": createItemsDetailedItemsDataForAdminEmail(order),
			"subtotal": order.GetFormattedSubTotal(),
			"tax_percent": order.GetFormattedTaxRate(),
			"tax": order.GetFormattedTaxAmount(),
			"total_amount": order.GetFormattedTotal(),
		},
		TemplateID:  send_email.QUICKBOOKS_FINAL_INVOICE_TEMPLATE_ID,
	}
	return emailData
}