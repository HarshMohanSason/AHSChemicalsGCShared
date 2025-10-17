package layout

import (
	"fmt"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	pdfutils "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

type CancellationSummary struct {
	TableValues        []canvas.TableRow
	TotalSummaryValues [][]string
	Date               string
	Time               string
	Total              string
}

var (
	CancellationSummaryTableHeaders        = []string{"ORDER ID", "CUSTOMER", "DATE", "ITEMS", "QTY", "PRICE PER UNIT", "TOTAL"}
	CancellationSummaryTableColWidths      = []float64{30, 50, 30, 92, 20, 30, 30}
	CancellationTotalSummaryTableHeaders   = []string{"CUSTOMER", "TAX RATE", "TAX AMOUNT", "SUBTOTAL", "TOTAL"}
	CancellationTotalSummaryTableColWidths = []float64{60, 20, 40, 40, 40}
)

func NewCancellationSummary(orders []*models.Order) *CancellationSummary {
	currentTime := time.Now().UTC()
	return &CancellationSummary{
		TableValues:        getTableValues(orders),
		TotalSummaryValues: getTotalSummaryValues(orders),
		Total:              getCancelledOrderTotal(orders),
		Date:               currentTime.Format("01/02/2006"),
		Time:               currentTime.Format("03:04:05 PM"),
	}
}

func getTableValues(orders []*models.Order) []canvas.TableRow {
	tableValues := make([]canvas.TableRow, 0)
	for _, order := range orders {

		tableCellValues := make([]canvas.TableCell, 0)

		tableCellValues = append(tableCellValues, canvas.TableCell{
			Lines: []string{order.ID},
			Width: CancellationSummaryTableColWidths[0],
		})
		tableCellValues = append(tableCellValues, canvas.TableCell{
			Lines: []string{order.Customer.GetFormattedName()},
			Width: CancellationSummaryTableColWidths[1],
		})
		tableCellValues = append(tableCellValues, canvas.TableCell{
			Lines: []string{order.CreatedAt.Format("01/02/2006")},
			Width: CancellationSummaryTableColWidths[2],
		})

		shortDescriptions := make([]string, 0)
		quantities := make([]string, 0)
		prices := make([]string, 0)
		totals := make([]string, 0)

		for _, item := range order.Items {
			shortDescriptions = append(shortDescriptions, item.GetShortDescription())
			quantities = append(quantities, item.GetFormattedQuantity())
			prices = append(prices, item.GetFormattedUnitPrice())
			totals = append(totals, item.GetFormattedTotalPrice())
		}

		tableCellValues = append(tableCellValues, canvas.TableCell{
			Lines: shortDescriptions,
			Width: CancellationSummaryTableColWidths[3],
		})
		tableCellValues = append(tableCellValues, canvas.TableCell{
			Lines: quantities,
			Width: CancellationSummaryTableColWidths[4],
		})
		tableCellValues = append(tableCellValues, canvas.TableCell{
			Lines: prices,
			Width: CancellationSummaryTableColWidths[5],
		})
		tableCellValues = append(tableCellValues, canvas.TableCell{
			Lines: totals,
			Width: CancellationSummaryTableColWidths[6],
		})
		tableValues = append(tableValues, canvas.TableRow{
			Cells: tableCellValues,
		})
	}
	return tableValues
}

func getTotalSummaryValues(orders []*models.Order) [][]string {
	totalSummaryValues := make([][]string, 0)
	for _, order := range orders {
		totalSummaryValues = append(totalSummaryValues, []string{
			order.Customer.Name,
			order.GetFormattedTaxRate(),
			order.GetFormattedTaxAmount(),
			order.GetFormattedSubTotal(),
			order.GetFormattedTotal(),
		})
	}
	return totalSummaryValues
}

func getCancelledOrderTotal(orders []*models.Order) string {
	total := 0.0
	for _, order := range orders {
		total += order.Total
	}
	return fmt.Sprintf("$%.2f", total)
}

func (cm *CancellationSummary) RenderToPDF() ([]byte, error) {

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pageWidth, pageHeight := pdf.GetPageSize()

	c := canvas.NewCanvas(pdf)
	c.SetBorderHeight(pageHeight - 6)
	c.SetBorderWidth(pageWidth - 6)
	c.SetBorderX(3)
	c.SetBorderY(3)
	c.SetMarginLeft(c.BorderX + 5)
	c.SetMarginTop(c.BorderY + 5)

	//Draw the outer border
	c.DrawRectangle(&canvas.Rectangle{
		X:           c.BorderX,
		Y:           c.BorderY,
		Width:       c.BorderWidth,
		Height:      c.BorderHeight,
		LineWidth:   0.8,
		BorderColor: canvas.PrimaryBlue,
	})
	c.IncY(5)
	c.IncX(5)

	c.DrawImageFromURL(canvas.ImageElement{
		URL:    company_details.LOGOPATH,
		X:      c.X,
		Y:      c.Y,
		Width:  60,
		Height: 0,
	})

	c.IncY(10)
	headingText := &canvas.Text{
		Content: "Cancellation Summary",
		Size:    16,
		X:       c.X,
		Y:       c.Y,
		Style:   "B",
		Font:    "Arial",
		Color:   canvas.PrimaryBlue,
	}
	c.MoveTo((pageWidth-headingText.GetWidth(c.PDF))/2, c.Y)
	headingText.SetX(c.X)
	c.DrawSingleLineText(headingText)

	c.MoveTo(c.BorderWidth-32, c.Y-5)
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Date:",
		Size:    10,
		X:       c.X,
		Y:       c.Y,
		Style:   "B",
		Font:    "Arial",
		Color:   canvas.Black,
	}, cm.Date)
	c.IncY(5)

	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Time:",
		Size:    10,
		X:       c.X,
		Y:       c.Y,
		Style:   "B",
		Font:    "Arial",
		Color:   canvas.Black,
	}, cm.Time)
	c.IncY(20)
	c.MoveTo(c.MarginLeft, c.Y)

	summaryTableHeader := &canvas.TableHeader{
		X:           c.X,
		Y:           c.Y,
		Headers:     CancellationSummaryTableHeaders,
		CellWidths:  CancellationSummaryTableColWidths,
		TextColor:   canvas.White,
		FillColor:   canvas.PrimaryBlue,
		BorderColor: canvas.PrimaryBlue,
	}
	summaryTableHeader.Draw(c, &canvas.Text{
		Size:    10,
		Style:   "B",
		Font:    "Arial",
		Color:   canvas.White,
	})
	summarTableBody := &canvas.TableBody2{
		X:           c.X,
		CellWidths:  CancellationSummaryTableColWidths,
		TextColor:   canvas.Black,
		BorderColor: canvas.PrimaryBlue,
		Rows: cm.TableValues,
	}
	summarTableBody.Y = summaryTableHeader.Y + summaryTableHeader.Height
	summarTableBody.Draw(c, &canvas.Text{
		Size:    9,
		Font:    "Arial",
		Color:   canvas.Black,
	})
	
	/*
	c.AddNewPageIfEnd(30, canvas.PrimaryBlue, 0.8)
	
	c.MoveTo(c.MarginLeft, tableEndYPos+10)
	c.DrawSingleLineText(&canvas.Text{
		Content: "Summary Section - Totals",
		Size:    11,
		X:       c.MarginLeft,
		Y:       c.Y,
		Style:   "B",
		Font:    "Arial",
		Color:   canvas.Black,
	})
	c.IncY(5)
	c.IncX(45)

	tableEndYPos = (&canvas.Table{
		Header: &canvas.TableHeader{
			X:           c.X,
			Y:           c.Y,
			Headers:     CancellationTotalSummaryTableHeaders,
			CellWidths:  CancellationTotalSummaryTableColWidths,
			TextColor:   canvas.White,
			FillColor:   canvas.PrimaryBlue,
			BorderColor: canvas.PrimaryBlue,
		},
		Body: &canvas.TableBody{
			X:           c.X,
			Y:           c.Y,
			CellWidths:  CancellationTotalSummaryTableColWidths,
			Rows:        cm.TotalSummaryValues,
			BorderColor: canvas.PrimaryBlue,
		},
		Width: pdfutils.CalculateTableCellWidths(CancellationTotalSummaryTableColWidths),
	}).Draw(c, &canvas.Text{
		Font:  "Helvetica",
		Size:  10,
		Style: "B",
	})
	c.MoveTo(c.X+119, tableEndYPos+5)

	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "TOTAL AMOUNT CANCELLED:",
		Size:    10,
		X:       c.X,
		Color:   canvas.Black,
		Y:       c.Y,
		Style:   "B",
		Font:    "Arial",
	}, cm.Total)
	*/
	c.DrawFooter("This is an automated document. Please do not reply to this email.")
	//Generate the PDF
	bytes, err := pdfutils.GetGeneratedPDF(pdf)
	return bytes, err
}