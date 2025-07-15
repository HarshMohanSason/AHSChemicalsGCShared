package canvas

import (
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/customers"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/orders"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/elements"
	"github.com/phpdave11/gofpdf"
)

const (
	TableWidth         = 179.0
	TableHeaderHeight  = 6.0
	TableHeaderPadding = 5.0 // Padding between the table header and the first row.
)

var (
	PrimaryBlue             = [3]int{65, 83, 145}
	PrimaryGreen            = [3]int{165, 199, 89}
	White                   = [3]int{255, 255, 255}
	Black                   = [3]int{0, 0, 0}
	ShippingTableHeaders    = []string{"REQUISITIONER", "SHIP VIA", "F.O.B", "SHIPPING TERMS"}
	ShippingTableValues     = [][]string{{"Robert Vodka", "In House", "Factory", "N/A"}}
	ProductTableHeaders     = []string{"SKU", "DESCRIPTION", "QTY", "PRICE", "TOTAL"}
	ShippingManifestHeaders = []string{"NUMBER OF UNITS", "HM TYPE", "CONTAINER DESCRIPTION AND CLASSIFICATION", "CLASS", "PRODUCT ID", "NET WEIGHT", "GROSS WEIGHT", "GROSS WEIGHT NHM", "GROSS WEIGHT HM"}
	InvoiceTableHeaders     = []string{"Item", "Quantity", "Price per unit", "Amount"}
)

type Canvas struct {
	PDF              *gofpdf.Fpdf
	X, Y             float64
	BorderX, BorderY float64 // Border starting point on page
	BorderWidth      float64
	BorderHeight     float64
	MarginLeft       float64 // Margin starting from the border
	MarginTop        float64
}

func NewCanvas(pdf *gofpdf.Fpdf) *Canvas {
	return &Canvas{
		PDF:          pdf,
		X:            0,
		Y:            0,
		BorderX:      8,
		BorderY:      8,
		BorderWidth:  193,
		BorderHeight: 280,
		MarginLeft:   15,
		MarginTop:    15,
	}
}

// Position helpers

func (c *Canvas) MoveTo(x, y float64) { c.X = x; c.Y = y }
func (c *Canvas) IncX(dx float64)     { c.X += dx }
func (c *Canvas) IncY(dy float64)     { c.Y += dy }
func (c *Canvas) DecX(dx float64)     { c.X -= dx }
func (c *Canvas) DecY(dy float64)     { c.Y -= dy }
func (c *Canvas) ResetX()             { c.X = c.BorderX + c.MarginLeft }
func (c *Canvas) ResetY()             { c.Y = c.BorderY + c.MarginTop }

// Draw border rectangle
func (c *Canvas) DrawBorder(width, height float64, thickness float64, color [3]int) {
	c.PDF.SetDrawColor(color[0], color[1], color[2])
	c.PDF.SetLineWidth(thickness)
	c.PDF.Rect(c.X, c.Y, width, height, "D")
}

// Draw company logo at current position
func (c *Canvas) DrawCompanyLogo(width, height float64) {
	imageElement := elements.Image{
		ImageURL: company_details.LOGOPATH,
		X:        c.X,
		Y:        c.Y,
		Width:    width,
		Height:   height,
		Flow:     false,
	}
	imageElement.Draw(c.PDF)
}

// Draw title text with alignment inside border
func (c *Canvas) DrawPDFTitle(title string, color [3]int, alignment string) {
	text := elements.Text{
		Text:  title,
		Font:  "Arial",
		Style: "B",
		X:     c.X,
		Y:     c.Y,
		Size:  24,
		Color: color,
	}
	text.ApplyTextStyle(c.PDF)

	textWidth := text.GetTextWidth(c.PDF)

	var xAlignment float64

	switch alignment {
	case "right":
		xAlignment = c.BorderWidth - textWidth
	case "left":
		xAlignment = c.MarginLeft
	default: // center
		xAlignment = (c.BorderWidth / 2) - (textWidth / 2) + c.MarginLeft
	}

	c.MoveTo(xAlignment, 33)
	text.X = c.X
	text.Y = c.Y
	text.Draw(c.PDF)
}

func (c *Canvas) DrawCompanyDetails() {
	lines := []struct {
		text  string
		style string
	}{
		{company_details.COMPANYNAME, "B"},
		{company_details.COMPANYADDRESSLINE1, ""},
		{company_details.COMPANYADDRESSLINE2, ""},
		{"Phone: " + company_details.COMPANYPHONE, ""},
		{"Email: " + company_details.COMPANYEMAIL, ""},
		{"Website: " + company_details.COMPANYWEBSITEURL, ""},
	}

	for _, line := range lines {
		textElement := elements.Text{
			Text:  line.text,
			Font:  "Arial",
			Style: line.style,
			Size:  10,
			X:     c.X,
			Y:     c.Y,
			Color: Black,
		}
		textElement.Draw(c.PDF)
		c.IncY(5)
	}
}

func (c *Canvas) DrawCustomerDetails(customer customers.Customer) {
	lines := []struct {
		text  string
		style string
	}{
		{customer.DisplayName, "B"},
		{customer.BillAddr.Line1, ""},
		{fmt.Sprintf("%s, %s %s", customer.BillAddr.City, customer.BillAddr.CountrySubDivisionCode, customer.BillAddr.PostalCode), ""},
		{"Phone: " + customer.PrimaryPhone.FreeFormNumber, ""},
		{"Email: " + customer.PrimaryEmailAddr.Address, ""},
	}

	for _, line := range lines {
		textElement := elements.Text{
			Text:  line.text,
			Font:  "Arial",
			Style: line.style,
			Size:  10,
			X:     c.X,
			Y:     c.Y,
			Color: Black,
		}
		textElement.Draw(c.PDF)
		c.IncY(5)
	}
}

func (c *Canvas) DrawTableHeaders(headers []string, colWidths []float64, fillColor [3]int, textColor [3]int) {

	for i, header := range headers {
		//Draw the rectangle
		rect := elements.Rectangle{
			X:         c.X,
			Y:         c.Y,
			Width:     colWidths[i],
			Height:    TableHeaderHeight,
			Style:     "F",
			FillColor: fillColor,
		}
		rect.Draw(c.PDF)
		c.IncY(TableHeaderHeight / 2)

		//Draw the text inside the rectangle
		textElement := elements.Text{
			Text:  header,
			Font:  "Arial",
			Style: "B",
			X:     c.X,
			Y:     c.Y + 0.5,
			Size:  10,
			Color: textColor,
		}
		textElement.DrawMultipleLines(c.PDF, colWidths[i], "center")

		//Reset positions to draw next cell
		c.DecY(TableHeaderHeight / 2)
		c.IncX(colWidths[i])
	}
}

func (c *Canvas) DrawTableRows(
	values [][]string,
	colWidths []float64,
	align string,
	fillColor, textColor [3]int,
	lineHeight float64,
) float64 {
	startYPos := c.Y
	tableHeight := 0.0

	for _, row := range values {
		// Calculate max height of each row before drawing it
		maxRowHeight := 0.0
		for j, col := range row {
			textElement := elements.Text{
				Text: col,
				X:    c.X,
				Y:    c.Y,
			}
			textHeight := textElement.GetMultiLineHeight(c.PDF, colWidths[j], lineHeight)
			maxRowHeight = max(maxRowHeight, textHeight) + 1.0
		}

		// Check if adding this row will overflow the page
		if c.Y+maxRowHeight > c.BorderHeight {
			//Finish drawing this table
			tableHeight += TableHeaderHeight + TableHeaderPadding
			c.MoveTo(c.MarginLeft, startYPos-TableHeaderHeight-TableHeaderPadding)
			c.DrawBorder(TableWidth, tableHeight, 0.8, PrimaryBlue)
			c.DrawTableCellRightBorder(len(values[0])-1, colWidths, 0.8, tableHeight, PrimaryBlue)

			// Start new page
			c.PDF.AddPage()
			c.MoveTo(c.BorderX, c.BorderY)
			c.DrawBorder(c.BorderWidth, c.BorderHeight, 0.8, PrimaryBlue)

			c.MoveTo(c.MarginLeft, c.MarginTop+15)
			startYPos = c.Y
			maxRowHeight = 0.0
			tableHeight = 0.0
		}

		// Now drawing the row safely
		for j, col := range row {
			textElement := elements.Text{
				Text: col,
				X:    c.X,
				Y:    c.Y,
			}
			textElement.DrawMultipleLines(c.PDF, colWidths[j], align)
			c.IncX(colWidths[j])
		}

		c.DecX(TableWidth)
		c.IncY(maxRowHeight)
		tableHeight += maxRowHeight
	}
	c.MoveTo(c.MarginLeft, startYPos-TableHeaderHeight-TableHeaderPadding) //Resetting the Y position to draw the entire table.
	return tableHeight + TableHeaderHeight
}

func (c *Canvas) DrawTableCellRightBorder(len int, colWidths []float64, thickness, tableHeight float64, borderColor [3]int) {

	for i := range len {
		xPos := c.X + colWidths[i]
		line := elements.Line{
			X1:    xPos,
			Y1:    c.Y,
			X2:    xPos,
			Y2:    c.Y + tableHeight,
			Width: thickness,
			Color: borderColor,
		}
		line.Draw(c.PDF)
		c.IncX(colWidths[i])
	}
}

func (c *Canvas) DrawBill(order *orders.Order) {
	const lineSpacing = 5.0
	const valueOffsetX = 35.0

	startX := c.X
	startY := c.Y

	labels := []string{
		"SUBTOTAL",
		fmt.Sprintf("TAX(%.2f%%)", order.TaxRate*100),
		"TOTAL",
	}
	values := []string{
		fmt.Sprintf("$%.2f", order.SubTotal),
		fmt.Sprintf("$%.2f", order.TaxAmount),
		fmt.Sprintf("$%.2f", order.Total),
	}

	for i := range labels {
		// Draw label
		label := elements.Text{
			Text:  labels[i],
			Font:  "Arial",
			Style: "",
			Size:  10,
			X:     startX,
			Y:     startY + float64(i)*lineSpacing,
			Color: Black,
		}
		label.Draw(c.PDF)

		// Draw value
		value := elements.Text{
			Text:  values[i],
			Font:  "Arial",
			Style: "",
			Size:  10,
			X:     startX + valueOffsetX,
			Y:     label.Y,
			Color: Black,
		}

		// Make TOTAL bold
		if i == len(labels)-1 {
			value.Style = "B"
		}
		value.Draw(c.PDF)
	}
}

func (c *Canvas) DrawFooter(pdfTitle string) {
	footerText := fmt.Sprintf(
		"If you have any questions or concerns about this %s please contact us at %s",
		pdfTitle, company_details.COMPANYEMAIL,
	)
	c.PDF.SetFont("Arial", "", 8)
	textWidth := c.PDF.GetStringWidth(footerText)
	textElement := elements.Text{
		Text:  footerText,
		X:     c.BorderX + ((c.BorderWidth / 2) - (textWidth / 2)),
		Y:     c.BorderY + c.BorderHeight - 5,
		Color: Black,
	}
	textElement.Draw(c.PDF)
}
