package invoice

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

func CreateInvoicePDF(order* models.Order, invoiceNo string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	c := canvas.NewCanvas(pdf)

	c.MoveTo(c.BorderX, c.BorderY)
	
	c.DrawBorder(c.BorderWidth, c.BorderHeight, 0.8, canvas.PrimaryGreen)
	c.MoveTo(c.MarginLeft, c.MarginTop + 10)
	
	c.DrawPDFTitle("INVOICE", canvas.PrimaryGreen, "left", 26)
	c.IncX(110)
	c.DecY(5)
	
	c.DrawCompanyLogo(65, 0)
	c.MoveTo(c.MarginLeft, c.MarginTop + 20)
	
	c.DrawCompanyDetails(10)
	c.IncY(5)
	
	customerStartXPos := c.X
	customerStartYPos := c.Y
	c.DrawText("Bill To ", "Arial", "B", 10)
	c.IncY(5)
	
	c.DrawCustomerDetails(&order.Customer)
	c.MoveTo(customerStartXPos + 125, customerStartYPos)
	
	c.DrawLabelWithText("Invoice No: ", invoiceNo, 10)
	c.IncY(5)
	
	c.DrawLabelWithText("Invoice Date: ", time.Now().Format("January 2, 2006"), 10)
	c.IncY(5)
	
	c.DrawLabelWithText("Payment Due: ", time.Now().AddDate(0, 0, 30).Format("January 2, 2006"), 10)
	c.IncY(25)
	c.ResetX()
	
	tableHeaderHeight := c.DrawTableHeaders(InvoiceTableHeaders, InvoiceTableColWidths, canvas.PrimaryGreen,canvas.White, 10)
	tableValues := utils.CreateTableRowValuesForInvoicePDF(order)
	
	var tableWidth float64 = 0.0
	for _, width := range InvoiceTableColWidths {
		tableWidth += width
	}

	//Set the fonts to get appropriate line height
	c.PDF.SetFont("Arial", "", 9)
	c.PDF.SetTextColor(canvas.Black[0], canvas.Black[1], canvas.Black[2])
	fontSize, _ := c.PDF.GetFontSize()
	lineHeight := c.PDF.PointConvert(fontSize)

	c.ResetX()
	c.IncY(tableHeaderHeight + canvas.TableHeaderPadding)
	tableHeight := c.DrawTableRows(tableValues, InvoiceTableColWidths, tableHeaderHeight, tableWidth, "center", canvas.PrimaryGreen, canvas.White, lineHeight)
	c.DrawBorder(tableWidth, tableHeight, 0.8, canvas.PrimaryGreen)
	c.DrawTableCellRightBorder(len(InvoiceTableHeaders) - 1, InvoiceTableColWidths, 0.8, tableHeight, canvas.PrimaryGreen)
	c.IncY(tableHeight + 5)
	c.MoveTo(128, c.Y)

	c.DrawBill(order)
	c.ResetX()
	c.IncY(20)

	c.DrawText("Terms & Conditions: ", "Arial", "B", 10)
	c.IncY(5)

	c.DrawMultiLineText(TermsAndConditions,  "Arial", "", 10, 100)
	
	c.DrawFooter("invoice")

	err := utils.GeneratePDFFileInPath(pdf, "invoice")
	if err != nil {
		//log.Print(err)
	}

	return "", err
}