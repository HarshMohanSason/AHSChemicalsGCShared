package purchase_order

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
)

func DrawPurchaseOrderShippingTable(c *canvas.Canvas) float64 {

	colWidths := []float64{35, 40, 35, 69}
	tableHeaderHeight := c.DrawTableHeaders(canvas.ShippingTableHeaders, colWidths, canvas.PrimaryBlue, canvas.White, 10)

	c.DecX(canvas.TableWidth)
	c.IncY(tableHeaderHeight + canvas.TableHeaderPadding)

	//Set the fonts to get appropriate line height
	c.PDF.SetFont("Arial", "", 9)
	c.PDF.SetTextColor(canvas.Black[0], canvas.Black[1], canvas.Black[2])
	fontSize, _ := c.PDF.GetFontSize()
	lineHeight := c.PDF.PointConvert(fontSize)

	tableHeight := c.DrawTableRows(canvas.ShippingTableValues, colWidths, tableHeaderHeight, 179, "center", canvas.PrimaryBlue, canvas.Black, lineHeight)
	c.DrawBorder(canvas.TableWidth, tableHeight, 0.8, canvas.PrimaryBlue)
	c.DrawTableCellRightBorder(len(canvas.ShippingTableHeaders)-1, colWidths, 0.8, tableHeight, canvas.PrimaryBlue)
	return c.Y + tableHeight
}

func DrawPurchaseOrderProductsTable(order *models.Order, c *canvas.Canvas) float64 {

	colWidths := []float64{30, 74, 15, 30, 30}
	tableHeaderHeight := c.DrawTableHeaders(canvas.ProductTableHeaders, colWidths, canvas.PrimaryBlue, canvas.White, 10)

	c.DecX(canvas.TableWidth)
	c.IncY(tableHeaderHeight + canvas.TableHeaderPadding)

	//Set the fonts to get appropriate line height
	c.PDF.SetFont("Arial", "", 9)
	c.PDF.SetTextColor(canvas.Black[0], canvas.Black[1], canvas.Black[2])
	fontSize, _ := c.PDF.GetFontSize()
	lineHeight := c.PDF.PointConvert(fontSize)

	mappedOrders := orders.CreateTableRowValuesForPurchaseOrderPDF(order)

	tableHeight := c.DrawTableRows(mappedOrders, colWidths, tableHeaderHeight, 179,"center", canvas.PrimaryBlue, canvas.Black, lineHeight)
	c.DrawBorder(canvas.TableWidth, tableHeight, 0.8, canvas.PrimaryBlue)
	c.DrawTableCellRightBorder(len(canvas.ProductTableHeaders)-1, colWidths, 0.8, tableHeight, canvas.PrimaryBlue)
	return c.Y + tableHeight
}
