package layout

import (
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

type OrderRevenueReport struct {
	InvoiceNo   string
	OrderNo     string
	Items       []*models.Product
	Customer    *models.Customer
	TableValues [][]string
	Cash        string
	COG         string // cost of goods
	TotalSales  string
	SalesTax    string
	TaxRate     string
	Revenue     string
	CreatedAt   string
}

var (
	orderRevenueReportHeaders        = []string{"SKU", "DESCRIPTION", "QTY", "SELLING PRICE", "PURCHASE PRICE", "TOTAL SELLING PRICE", "TOTAL PURCHASE PRICE", "TOTAL REVENUE"}
	orderRevenueReportTableColWidths = []float64{25, 35, 10, 22, 22, 22, 22, 22}
)

func NewOrderRevenueReport(order *models.Order, invoiceNo string) *OrderRevenueReport {
	orderRevenue := &OrderRevenueReport{
		InvoiceNo:  invoiceNo,
		OrderNo:    order.ID,
		Items:      order.Items,
		Customer:   order.Customer,
		Cash:       order.GetFormattedTotal(),
		COG:        order.GetFormattedCOG(),
		TotalSales: order.GetFormattedSubTotal(),
		SalesTax:   order.GetFormattedTaxAmount(),
		TaxRate:    order.GetFormattedTaxRate(),
		Revenue:    order.GetFormattedTotalRevenue(),
		CreatedAt:  order.UpdatedAt.Format("January 2, 2006"),
	}
	orderRevenue.CreateTableValues(order)
	return orderRevenue
}

func (orr *OrderRevenueReport) CreateTableValues(order *models.Order) {
	tableValues := make([][]string, 0)
	for _, item := range order.Items {
		tableValues = append(tableValues, []string{item.SKU, item.GetFormattedDescription(), item.GetFormattedQuantity(), item.GetFormattedUnitPrice(), item.GetFormattedPurchasePrice(), item.GetFormattedTotalPrice(), item.GetFormattedTotalPurchasePrice(), item.GetFormattedTotalRevenue()})
	}
	orr.TableValues = tableValues
}

func (orr *OrderRevenueReport) RenderToPDF() ([]byte, error) {

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

	//Draw the company logo on top left
	c.DrawSingleLineText(&canvas.Text{
		Content: "REVENUE REPORT",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    24,
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
		Content: "Customer Details",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	})
	c.IncY(5)

	c.DrawCustomerDetails(orr.Customer)
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
	}, orr.InvoiceNo)
	c.IncY(5)

	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Order No:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, orr.OrderNo)
	c.IncY(5)

	//Invoice Date
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Created At:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, orr.CreatedAt)
	c.IncY(5)

	c.IncY(25)
	c.ResetX()

	//Draw the table
	tableEndYPos := (&canvas.Table{
		Header: &canvas.TableHeader{
			X:           c.X,
			Y:           c.Y,
			Headers:     orderRevenueReportHeaders,
			CellWidths:  orderRevenueReportTableColWidths,
			TextColor:   canvas.White,
			FillColor:   canvas.PrimaryGreen,
			BorderColor: canvas.PrimaryGreen,
		},
		Body: &canvas.TableBody{
			X:           c.X,
			Y:           c.Y,
			CellWidths:  orderRevenueReportTableColWidths,
			Rows:        orr.TableValues,
			TextColor:   canvas.Black,
			BorderColor: canvas.PrimaryGreen,
		},
		Width: utils.CalculateShippingTableCellWidths(orderRevenueReportTableColWidths),
	}).Draw(c, &canvas.Text{
		Font:  "Helvetica",
		Size:  10,
		Style: "B",
		Color: canvas.White,
	})
	c.MoveTo(c.MarginLeft, tableEndYPos+5)

	//Check if a new page is needed or not
	c.AddNewPageIfEnd(10, canvas.PrimaryGreen, 0.8)

	c.IncX(108)
	c.DrawBillingDetails([]string{"CASH", "COST OF GOODS", "SALES", fmt.Sprintf("SALES TAX AMOUNT (%s)", orr.TaxRate), "TOTAL REVENUE"}, []string{orr.Cash, orr.COG, orr.TotalSales, orr.SalesTax, orr.Revenue}, false, false)
	
	//Generate the PDF
	bytes, err := utils.GetGeneratedPDF(c.PDF)
	return bytes, err
}
