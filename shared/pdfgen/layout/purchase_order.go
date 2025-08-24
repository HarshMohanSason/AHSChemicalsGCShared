package layout

import (
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

var (
	shippingTableHeaders    = []string{"REQUISITIONER", "SHIP VIA", "F.O.B", "SHIPPING TERMS"}
	shippingTableValues     = [][]string{{"N/A", "In House", "Factory", "N/A"}}
	shippingTableCellWidths = []float64{60, 40, 30, 50}
	productTableHeaders     = []string{"SKU", "DESCRIPTION", "QTY", "PRICE", "TOTAL"}
	productTableCellWidths  = []float64{30, 65, 25, 30, 30}
)

type PurchaseOrder struct {
	ID                  string
	Customer            models.Customer
	SpecialInstructions string
	TableValues         [][]string
	TaxRate             string
	TaxAmount           string
	SubTotal            string
	Total               string
	CreatedAt           string
}

func NewPurchaseOrder(o *models.Order) *PurchaseOrder {
	purchaseOrder := &PurchaseOrder{
		ID:                  o.ID,
		Customer:            o.Customer,
		SpecialInstructions: o.SpecialInstructions,
		TaxRate:             o.GetFormattedTaxRate(),
		TaxAmount:           o.GetFormattedTaxAmount(),
		SubTotal:            o.GetFormattedSubTotal(),
		Total:               o.GetFormattedTotal(),
		CreatedAt:           o.CreatedAt.Format("January 2, 2006 at 3:04 PM"),
	}
	purchaseOrder.getTableValues(o.Items)
	return purchaseOrder
}

func (p *PurchaseOrder) getTableValues(items []models.Product) {
	tableValues := make([][]string, 0)
	for _, item := range items {
		tableValues = append(tableValues, []string{
			item.SKU,
			item.GetFormattedDescription(),
			item.GetFormattedQuantity(),
			item.GetFormattedUnitPrice(),
			item.GetFormattedTotalPrice(),
		})
	}
	p.TableValues = tableValues
}

func (p *PurchaseOrder) RenderToPDF() ([]byte, error) {
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
		BorderColor: canvas.PrimaryBlue,
	})
	c.MoveTo(c.MarginLeft, 10)

	//Draw the company logo on top left
	c.DrawImageFromURL(canvas.ImageElement{
		URL:    company_details.LOGOPATH,
		X:      c.X,
		Y:      c.Y,
		Width:  65,
		Height: 0,
	})
	c.IncX(100)
	c.IncY(25)

	//Draw the PDF Name on top right side
	c.DrawSingleLineText(&canvas.Text{
		Content: "PURCHASE ORDER",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    24,
		Color:   canvas.PrimaryBlue,
		Style:   "B",
	})
	c.ResetX()
	c.IncY(7)

	c.DrawCompanyDetails()
	c.IncX(120)
	c.DecY(30)

	//Draw the Date
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Date:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.CreatedAt)
	c.IncY(5)

	//Draw the PO Number
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "P.O. #:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.ID)
	c.ResetX()
	c.IncY(25)

	//Draw the Vendor Heading
	c.DrawTextInColoredRect(&canvas.Text{
		Content: "VENDOR",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.White,
		Style:   "B",
	}, &canvas.Rectangle{
		X:         c.X,
		Y:         c.Y,
		Width:     75,
		Style:     "F",
		Height:    6,
		LineWidth: 0.8,
		FillColor: canvas.PrimaryBlue,
	}, "left")
	c.IncY(11)
	//Vendor Details
	c.DrawCompanyDetails()
	c.DecY(41)
	c.IncX(c.X + 90)

	//Draw Ship To Heading
	c.DrawTextInColoredRect(&canvas.Text{
		Content: "SHIP TO",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.White,
		Style:   "B",
	}, &canvas.Rectangle{
		X:         c.X,
		Y:         c.Y,
		Width:     75,
		Style:     "F",
		Height:    6,
		LineWidth: 0.8,
		FillColor: canvas.PrimaryBlue,
	}, "left")
	c.IncY(11)
	//Ship To Details
	c.DrawCustomerDetails(&p.Customer)
	c.IncY(5)
	c.ResetX()

	//Draw the shipping table
	tableEndYPos := (&canvas.Table{
		Header: &canvas.TableHeader{
			X:           c.X,
			Y:           c.Y,
			Headers:     shippingTableHeaders,
			CellWidths:  shippingTableCellWidths,
			TextColor:   canvas.White,
			FillColor:   canvas.PrimaryBlue,
			BorderColor: canvas.PrimaryBlue,
		},
		Body: &canvas.TableBody{
			X:           c.X,
			Y:           c.Y,
			CellWidths:  shippingTableCellWidths,
			Rows:        shippingTableValues,
			TextColor:   canvas.Black,
			BorderColor: canvas.PrimaryBlue,
		},
		Width: utils.CalculateShippingTableCellWidths(shippingTableCellWidths),
	}).Draw(c, &canvas.Text{
		Font:  "Helvetica",
		Size:  10,
		Style: "B",
		Color: canvas.White,
	})
	c.MoveTo(c.MarginLeft, tableEndYPos)
	c.IncY(5)

	//Draw the product table
	tableEndYPos = (&canvas.Table{
		Header: &canvas.TableHeader{
			X:           c.X,
			Y:           c.Y,
			Headers:     productTableHeaders,
			CellWidths:  productTableCellWidths,
			TextColor:   canvas.White,
			FillColor:   canvas.PrimaryBlue,
			BorderColor: canvas.PrimaryBlue,
		},
		Body: &canvas.TableBody{
			X:           c.X,
			Y:           c.Y,
			Height:      0,
			CellWidths:  productTableCellWidths,
			Rows:        p.TableValues,
			TextColor:   canvas.Black,
			BorderColor: canvas.PrimaryBlue,
		},
		Width: utils.CalculateShippingTableCellWidths(productTableCellWidths),
	}).Draw(c, &canvas.Text{
		Font:  "Helvetica",
		Size:  10,
		Style: "B",
	})
	c.MoveTo(147, tableEndYPos+5)

	//Check if new page needs to be drawn before drawing the billing section (30px is the height of the label drawn)
	c.AddNewPageIfEnd(30, canvas.PrimaryBlue, 0.8)

	//Draw the billing section
	c.DrawBillingDetails([]string{"SUBTOTAL", fmt.Sprintf("TAX (%s)", p.TaxRate), "TOTAL"}, []string{p.SubTotal, p.TaxAmount, p.Total}, false, false)
	c.MoveTo(c.MarginLeft, tableEndYPos+5)

	//Comments or Special Instructions label
	c.DrawSingleLineText(&canvas.Text{
		Content: "Comments or Special Instructions:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Style:   "B",
	})
	c.IncY(5)
	//Comments or Special Instructions
	c.DrawMultipleLines(&canvas.Text{
		Content: p.SpecialInstructions,
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    9,
		Style:   "",
	}, 90, "")

	//Footer
	c.DrawFooter(fmt.Sprintf("If you have any questions or concerns about this purchase order please contact us at %s", company_details.COMPANYEMAIL))

	//Generate the PDF
	bytes, err := utils.GetGeneratedPDF(c.PDF)
	return bytes, err
}
