package invoice

//Invoice constants vars for pdf generation
var (
	InvoiceTableHeaders   = []string{"ITEM", "QUANTITY", "PRICE PER UNIT", "AMOUNT"}
	InvoiceTableColWidths = []float64{60, 25, 45, 45}
	TermsAndConditions    = "Payment is due within 30 days from the invoice date (Net 30). A 10% late fee will be automatically applied to the total outstanding balance if payment is not received within 14 days after the due date. Continued non-payment may result in suspension of services and additional collection actions. By receiving this invoice, you agree to these terms."
)
