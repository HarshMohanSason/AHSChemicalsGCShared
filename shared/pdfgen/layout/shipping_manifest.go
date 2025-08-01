package layout

import (
	"bytes"
	"fmt"
	"image"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/company_details"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/utils"
	"github.com/phpdave11/gofpdf"
)

var (
	shippingManifestHeaders        = []string{"UNITS", "HM", "TYPE CONTAINER", "DESCRIPTION AND CLASSIFICATION", "CLASS", "SKU", "NET WEIGHT", "GROSS WEIGHT NHM", "GROSS WEIGHT HM"}
	shippingManifestTableColWidths = []float64{13, 10, 21, 40, 14, 30, 20.6, 20.6, 20.6}
	typeContainer                  = "Carton"
	class                          = "55.0"
)

type ShippingManifest struct {
	PONumber             string
	Customer             *models.Customer
	Product              []models.Product
	DeliveredBy          string
	TotalUnits           string
	TotalNonHazardWeight string
	TotalHazardousWeight string
	TotalWeight          string
	ReceivedBy           string
	Signature            []byte
	DeliverImages        [][]byte
	TableValues          [][]string
	DeliveredAt          string
}

func NewShippingManifest(delivery *models.Delivery) *ShippingManifest {
	if delivery == nil {
		return nil
	}
	shippingManifest := &ShippingManifest{
		PONumber:             delivery.Order.ID,
		Customer:             &delivery.Order.Customer,
		Product:              delivery.Order.Items,
		TotalUnits:           delivery.Order.GetFormattedTotalItems(),
		TotalNonHazardWeight: delivery.Order.GetFormattedNetNonHazardousWeight(),
		TotalHazardousWeight: delivery.Order.GetFormattedNetHazardousWeight(),
		TotalWeight:          delivery.Order.GetFormattedNetWeight(),
		DeliveredBy:          delivery.DeliveredBy,
		ReceivedBy:           delivery.ReceivedBy,
		Signature:            delivery.Signature,
		DeliverImages:        delivery.DeliveryImages,
		DeliveredAt:          delivery.DeliveredAt.Format("January 2, 2006 at 3:04 PM"),
	}
	shippingManifest.getTableValues(shippingManifest.Product)
	return shippingManifest
}

func (p *ShippingManifest) getTableValues(items []models.Product) {
	tableValues := make([][]string, 0)
	for _, item := range items {
		tableValues = append(tableValues, []string{
			item.GetFormattedQuantity(),
			item.GetFormattedIsHazardous(),
			typeContainer,
			item.GetFormattedDescription(),
			class,
			item.SKU,
			item.GetFormattedTotalWeight(),
			item.GetFormattedTotalHazardousWeight(),
			item.GetFormattedTotalNonHazardousWeight(),
		})
	}
	p.TableValues = tableValues
}

func (p *ShippingManifest) RenderToPDF() ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	c := canvas.NewCanvas(pdf)
	//Shipping manifest is has more table rows so it needs more space
	c.BorderX = 5
	c.BorderY = 5
	c.BorderWidth = 200
	c.BorderHeight = 285
	c.MarginLeft = c.BorderX + 5
	c.MarginTop = c.BorderX + 5

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
	c.IncX(105)
	c.IncY(25)

	//Draw the PDF Name on top right side
	c.DrawSingleLineText(&canvas.Text{
		Content: "SHIPPING MANIFEST",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    24,
		Color:   canvas.PrimaryBlue,
		Style:   "B",
	})
	c.ResetX()
	c.IncY(7)

	//Ship To Section
	c.DrawSingleLineText(&canvas.Text{
		Content: "Ship To",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	})
	c.IncY(5)

	c.DrawCustomerDetails(p.Customer)
	c.DecY(30)
	c.MoveTo(c.MarginLeft+105, c.Y)

	c.DrawCompanyDetails()
	c.IncY(5)
	c.ResetX()

	//24 Hour Contact Number
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "24 Hour Contact Number:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, company_details.COMPANY24HOURPHONE)
	c.MoveTo(c.MarginLeft+105, c.Y)

	//Delivered At Section
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Delivered At:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.DeliveredAt)
	c.IncY(5)

	//P.O. Number
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "P.O.#:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.PONumber)
	c.IncY(10)
	c.ResetX()

	tableEndYPos := (&canvas.Table{
		Header: &canvas.TableHeader{
			X:           c.X,
			Y:           c.Y,
			Headers:     shippingManifestHeaders,
			CellWidths:  shippingManifestTableColWidths,
			FillColor:   canvas.PrimaryBlue,
			BorderColor: canvas.PrimaryBlue,
			TextColor:   canvas.White,
		},
		Body: &canvas.TableBody{
			X:           c.X,
			Y:           c.Y,
			Rows:        p.TableValues,
			CellWidths:  shippingManifestTableColWidths,
			TextColor:   canvas.Black,
			BorderColor: canvas.PrimaryBlue,
		},
	}).Draw(c, &canvas.Text{
		Font:  "Helvetica",
		X:     c.X,
		Y:     c.Y,
		Size:  9,
		Color: canvas.White,
		Style: "B",
	})
	c.MoveTo(c.MarginLeft, tableEndYPos+5)

	//Total Units
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "Total Units: ",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.TotalUnits)
	c.IncX(120)

	//Check if a new page needs to be created
	c.AddNewPageIfEnd(10, canvas.PrimaryBlue, 0.8)

	//Total Non Hazardous Weight
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "NON HAZARDOUS WEIGHT: ",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.TotalNonHazardWeight)
	c.IncY(5)

	//Total Hazardous Weight
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "HAZARDOUS WEIGHT: ",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.TotalHazardousWeight)
	c.IncY(5)

	//Total Net Weight
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "TOTAL WEIGHT: ",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.TotalWeight)
	c.MoveTo(c.MarginLeft, c.Y+10)

	//Received By
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "RECEIVED BY: ",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.ReceivedBy)
	c.IncX(120)

	//Delivered By
	c.DrawLabelWithSingleLineText(&canvas.Text{
		Content: "DELIVERED BY:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	}, p.DeliveredBy)
	c.MoveTo(c.MarginLeft, c.Y+5)

	//Signature label
	c.DrawSingleLineText(&canvas.Text{
		Content: "SIGNATURE:",
		Font:    "Helvetica",
		X:       c.X,
		Y:       c.Y,
		Size:    10,
		Color:   canvas.Black,
		Style:   "B",
	})
	c.IncY(5)

	//Signature Image
	c.DrawImageFromBytes(canvas.ImageElement{
		X:      c.X,
		Y:      c.Y,
		Width:  35,
		Height: 0,
		Bytes:  p.Signature,
	})

	c.DrawFooter(fmt.Sprintf("If you have any questions or concerns about this shipping manifest please contact us at %s", company_details.COMPANYEMAIL))
	
	//Delivery Images
	for _, imageBytes := range p.DeliverImages {
		//Each image is drawn on a new page
		c.PDF.AddPage()
		c.MoveTo(0, 0)

		var imageWidth float64
		var imageHeight float64 
		pageWidth, pageHeight := c.PDF.GetPageSize()
		image, _, err := image.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			imageWidth = c.BorderWidth + c.BorderX //Page width basically 
			imageHeight = 0
		}else{
			imageWidth, imageHeight = c.GetCorrectByteImageDimensions(image)
		}

		//Delivery Image
		c.DrawImageFromBytes(canvas.ImageElement{
			X:      (pageWidth - imageWidth)/2, //Center horizontally
			Y:      (pageHeight - imageHeight)/2, //Center vertically
			Width:  imageWidth,
			Height: imageHeight,
			Bytes:  imageBytes,
		})
	}

	//Generate the PDF
	bytes, err := utils.GetGeneratedPDF(c.PDF)
	return bytes, err
}
