package purchase_order

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/orders"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

func CreatePurchaseOrderPDF(order *orders.Order) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	c := canvas.NewCanvas(pdf)
	c.MoveTo(c.BorderX, c.BorderY)
	c.DrawBorder(c.BorderWidth, c.BorderHeight, 0.8, canvas.PrimaryBlue)

	c.MoveTo(c.MarginLeft, 10)
	c.DrawCompanyLogo(65, 0)

	c.DrawPDFTitle("PURCHASE ORDER", canvas.PrimaryBlue, "right")

	c.MoveTo(c.MarginLeft, 40)

	c.DrawCompanyDetails()

	c.MoveTo(c.BorderWidth-45, 40)

	DrawDateAndPONumber(order, c)

	c.MoveTo(c.MarginLeft, 70)

	DrawVendorShipToSectionHeading("VENDOR", c)
	c.IncY(10)
	c.DrawCompanyDetails()

	c.MoveTo(c.BorderWidth-70, 70)

	DrawVendorShipToSectionHeading("SHIP TO", c)
	c.IncY(10)
	c.DrawCustomerDetails(order.Customer)

	c.MoveTo(c.MarginLeft, 110)
	tablePos := DrawPurchaseOrderShippingTable(c) //Track the y position of the table to draw the next item

	c.MoveTo(c.MarginLeft, tablePos+5)
	tablePos = DrawPurchaseOrderProductsTable(order, c)

	c.MoveTo(c.MarginLeft, tablePos+10)
	DrawCommentsorSpecialInstructions(order.SpecialInstructions, c)

	c.MoveTo(c.BorderWidth - 58, tablePos+5)
	c.DrawBill(order)

	c.DrawFooter("purchase order")
	if err := utils.GeneratePDFFileInPath(pdf, "purchase_order"); err != nil {
		panic(err)
	}
	return pdf
}
