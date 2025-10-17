// package layout contains the layout of any pdf generated. Each layout contains a function RenderToPDF which returns the appropriate structure with a function abiding to the interface defined in package pdfgen
package layout

import (
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	pdfutils "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

type Invoice struct {
	Number      string
	Items       []*models.Product
	Customer    *models.Customer
	TableValues [][]string
	LateFee     string
	Total       string
	SubTotal    string
	TaxAmount   string
	TaxRate     string
	PaymentDue  string
	LateFeeDate string
	CreatedAt   string
}

const (
	TermsAndConditions = "The payment for this invoice is due within 30 days from the invoice date (Net 30). By receiving this invoice, you agree to these terms."
)

var (
	invoiceTableHeaders   = []string{"ITEM", "QUANTITY", "PRICE PER UNIT", "AMOUNT"}
	invoiceTableColWidths = []float64{75, 25, 40, 40}
)

func NewInvoice(order *models.Order, invoiceNumber string) *Invoice {
	invoice := &Invoice{
		Number:      invoiceNumber,
		Customer:    order.Customer,
		LateFee:     fmt.Sprintf("$%.2f", order.Total*0.1),
		Total:       order.GetFormattedTotal(),
		SubTotal:    order.GetFormattedSubTotal(),
		TaxAmount:   order.GetFormattedTaxAmount(),
		TaxRate:     order.GetFormattedTaxRate(),
		CreatedAt:   order.GetLocalUpdatedAtTime().Format("January 2, 2006"),
		PaymentDue:  order.GetLocalUpdatedAtTime().AddDate(0, 0, 30).Format("January 2, 2006"),
		LateFeeDate: order.GetLocalUpdatedAtTime().AddDate(0, 0, 44).Format("January 2, 2006"),
	}
	invoice.setTableValues(order.Items)

	return invoice
}

func (i *Invoice) setTableValues(items []*models.Product) {
	tableValues := make([][]string, 0)
	for _, item := range items {
		tableValues = append(tableValues, []string{
			item.GetFormattedDescription(),
			item.GetFormattedQuantity(),
			item.GetFormattedUnitPrice(),
			item.GetFormattedTotalPrice(),
		})
	}
	i.TableValues = tableValues
}

func (i *Invoice) RenderToPDF() ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	c := canvas.NewCanvas(pdf)
	c.MoveTo(c.BorderX, c.BorderY)

	//Draw the outer border
	c.DrawRectangle(&canvas.Rectangle{
		X:           c.BorderX,
		Y:           c.BorderY,
		Width:       c.BorderWidth,
		Height:      c.BorderHeight,
		LineWidth:   0.8,
		BorderColor: canvas.PrimaryGreen,
	})
	c.MoveTo(c.MarginLeft, c.MarginTop+10)

	//Draw the pdf title on top left
	c.DrawSingleLineText(&canvas.Text{
		Content: "INVOICE",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    26,
		Color:   canvas.PrimaryGreen,
		Style:   "B",
	})
	c.IncY(10)

	c.DrawCompanyDetails()
	companyDetailsEndYPos := c.Y
	c.MoveTo(125, c.MarginTop)

	//Draw the company logo on top right
	c.DrawImageFromURL(canvas.ImageElement{
		URL:    company_details.LOGOPATH,
		X:      c.X,
		Y:      c.Y,
		Width:  70,
		Height: 0,
	})
	c.MoveTo(c.MarginLeft, companyDetailsEndYPos+5)

	c.DrawSingleLineText(&canvas.Text{
		Content: "Bill To",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	})
	c.IncY(5)

	c.DrawCustomerDetails(i.Customer)
	c.MoveTo(135, companyDetailsEndYPos+5)

	//Invoice No
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Invoice No:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, i.Number)
	c.IncY(5)

	//Invoice Date
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Invoice Date:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, i.CreatedAt)
	c.IncY(5)

	//Payment Due
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Payment Due:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, i.PaymentDue)
	c.IncY(5)

	//Late Fee Due
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Late Fee Date:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, i.LateFeeDate)

	c.IncY(25)
	c.ResetX()

	tableEndYPos := (&canvas.Table{
		Header: &canvas.TableHeader{
			X:           c.X,
			Y:           c.Y,
			Headers:     invoiceTableHeaders,
			CellWidths:  invoiceTableColWidths,
			TextColor:   canvas.White,
			FillColor:   canvas.PrimaryGreen,
			BorderColor: canvas.PrimaryGreen,
			BorderThickness: 0.8,
		},
		Body: &canvas.TableBody{
			X:           c.X,
			Y:           c.Y,
			CellWidths:  invoiceTableColWidths,
			Rows:        i.TableValues,
			TextColor:   canvas.Black,
			BorderColor: canvas.PrimaryGreen,
			BorderThickness: 0.8,
		},
		Width: pdfutils.CalculateTableCellWidths(shippingTableCellWidths),
	}).Draw(c, &canvas.Text{
		Font:  "Helvetica",
		Size:  10,
		Style: "B",
		Color: canvas.White,
	})
	c.MoveTo(c.MarginLeft, tableEndYPos+5)

	//Check if we need to add a new page (3 labels with a gap of 10px each, so 30px total)
	c.AddNewPageIfEnd(30, canvas.PrimaryGreen, 0.8)

	c.IncX(127)
	c.DrawBillingDetails([]string{"SUBTOTAL", fmt.Sprintf("TAX (%s)", i.TaxRate), "TOTAL"}, []string{i.SubTotal, i.TaxAmount, i.Total}, false, false)
	c.MoveTo(c.MarginLeft, c.Y+5)

	//Check if we need to add a new page (3 labels with a gap of 10px each, so 30px total)
	c.AddNewPageIfEnd(25, canvas.PrimaryGreen, 0.8)
	c.DrawSingleLineText(&canvas.Text{
		Content: "Notes",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	})
	c.IncY(5)

	//Check if we need to add a new page (3 labels with a gap of 10px each, so 30px total)
	c.AddNewPageIfEnd(25, canvas.PrimaryGreen, 0.8)
	c.DrawMultipleLines(&canvas.Text{
		Content: fmt.Sprintf("A late fee of %s will be charged if the invoice is not paid within the late fee date. Continued non-payment may result in suspension of services and additional collection actions.", i.LateFee),
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "",
	}, 100, "")
	c.IncY(15)

	c.DrawSingleLineText(&canvas.Text{
		Content: "Terms & Conditions",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	})
	c.IncY(5)
	c.DrawMultipleLines(&canvas.Text{
		Content: TermsAndConditions,
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "",
	}, 100, "")

	c.DrawFooter(fmt.Sprintf("If you have any questions or concerns about this invoice please contact us at %s", company_details.COMPANYEMAIL))

	//Generate the PDF
	bytes, err := pdfutils.GetGeneratedPDF(c.PDF)
	return bytes, err
}
