package shipping_manifest

import (
	"log"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

// order *orders.Order, deliveryData *orders.DeliveryData, signatureFile multipart.File, deliveryImages []multipart.File
func CreateShippingManifestPDF(order *models.Order, deliveryData *models.Delivery) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	c := canvas.NewCanvas(pdf)
	//Setting custom values since shipping manifest has a lot more table headers than other PDFs
	c.BorderX = 5
	c.BorderY = 5
	c.MarginLeft = c.BorderX + 5
	c.MarginTop = c.BorderY + 5
	pageWidth, pageHeight := pdf.GetPageSize()
	c.BorderWidth = pageWidth - 10
	c.BorderHeight = pageHeight - 10
	c.MoveTo(c.BorderX, c.BorderY)

	c.DrawBorder(c.BorderWidth, c.BorderHeight, 0.8, canvas.PrimaryBlue)
	c.MoveTo(c.MarginLeft, c.MarginTop)
	c.DrawCompanyLogo(65, 0)

	c.MoveTo(c.MarginLeft, 26)
	c.DrawPDFTitle("SHIPPING MANIFEST", canvas.PrimaryBlue, "right", 20)
	rightSideXPos := c.X
	c.IncY(8)

	c.DrawCompanyDetails(10)
	c.MoveTo(c.MarginLeft, c.Y-30)

	c.DrawText("Ship To", "Arial", "B", 10)
	c.IncY(5)

	c.DrawCustomerDetails(&order.Customer)
	c.IncY(3)

	c.DrawLabelWithText("24 Hour Contact Number: ", "+1 (800) 577-6202", 10)
	c.MoveTo(rightSideXPos, c.Y)

	c.DrawLabelWithText("Delivered at: ", time.Now().Format("January 2, 2006 at 03:04 PM"), 10)
	c.IncY(5)

	c.DrawLabelWithText("PO #: ", order.ID, 10)
	c.IncY(10)

	c.MoveTo(c.MarginLeft, c.Y)
	tableHeaderHeight := c.DrawTableHeaders(ShippingManifestHeaders, ShippingManifestTableColWidths, canvas.PrimaryBlue, canvas.White, 8)

	c.MoveTo(c.MarginLeft, c.Y+tableHeaderHeight)
	mappedTableValues := utils.CreateTableRowValuesForShippingManifestPDF(order)

	//Set the fonts to get appropriate line height
	c.PDF.SetFont("Arial", "", 8)
	c.PDF.SetTextColor(canvas.Black[0], canvas.Black[1], canvas.Black[2])
	fontSize, _ := c.PDF.GetFontSize()
	lineHeight := c.PDF.PointConvert(fontSize)
	c.IncY(canvas.TableHeaderPadding)

	var tableWidth float64 = 0.0
	for _, item := range ShippingManifestTableColWidths {
		tableWidth += item
	}
	tableHeight := c.DrawTableRows(mappedTableValues, ShippingManifestTableColWidths, tableHeaderHeight, tableWidth, "center", canvas.PrimaryBlue, canvas.White, lineHeight)

	c.DrawBorder(tableWidth, tableHeight, 0.8, canvas.PrimaryBlue)
	c.DrawTableCellRightBorder(len(ShippingManifestHeaders)-1, ShippingManifestTableColWidths, 0.8, tableHeight, canvas.PrimaryBlue)
	c.IncY(tableHeight + 10)
	c.ResetX()

	c.DrawLabelWithText("TOTAL UNITS: ", order.GetFormattedTotalItems(), 10)

	c.MoveTo(125, c.Y)
	c.DrawLabelWithText("NON HAZARDOUS WEIGHT: ", order.GetFormattedHazardousWeight(), 10)
	c.IncY(5)
	c.DrawLabelWithText("HAZARDOUS WEIGHT: ", order.GetFormattedHazardousWeight(), 10)
	c.IncY(5)
	c.DrawLabelWithText("TOTAL WEIGHT: ", order.GetFormattedNetWeight(), 10)

	c.IncY(10)
	c.ResetX()

	c.DrawLabelWithText("RECEIVED BY: ", deliveryData.ReceivedBy, 10)
	c.MoveTo(125, c.Y)
	c.DrawLabelWithText("DELIVERED BY: ", deliveryData.DeliveredBy, 10)
	c.IncY(5)
	c.ResetX()
	c.DrawText("SIGNATURE: ", "Arial", "B", 10)

	c.DrawFooter("shipping manifest")
	c.DrawImageFromMultiPart(deliveryData.Signature, 40, 0)

	//Draw the delivery images
	for _, image := range deliveryData.Images {
		c.PDF.AddPage()
		c.MoveTo(c.BorderX, c.BorderY)
		c.DrawBorder(c.BorderWidth, c.BorderHeight, 0.8, canvas.PrimaryBlue)
		c.MoveTo(c.MarginLeft, c.MarginTop)
		c.DrawImageFromMultiPart(image, c.BorderWidth - 10, 0)
	}

	err := utils.GeneratePDFFileInPath(pdf, "shipping_manifest")
	if err != nil {
		log.Print(err)
	}

	return "", err
}