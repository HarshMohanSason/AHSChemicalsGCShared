package purchase_order

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/orders"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/elements"
)

// DrawDateAndPONumber draws the date and purchase order number on the top right corner below the pdf title
func DrawDateAndPONumber(o *orders.Order, c *canvas.Canvas) {

	//Date
	label := elements.Text{
		Text:  "Date: ",
		Font:  "Arial",
		Style: "B",
		Size:  10,
		X:     c.X,
		Y:     c.Y,
		Color: canvas.Black,
	}
	label.Draw(c.PDF)

	value := elements.Text{
		Text:  o.CreatedAt.Format("January 2, 2006"),
		Font:  "Arial",
		Style: "",
		Size:  10,
		X:     c.X + label.GetTextWidth(c.PDF),
		Y:     c.Y,
		Color: canvas.Black,
	}
	value.Draw(c.PDF)

	c.IncY(5)

	// P.O. #
	label.Text = "P.O. #: "
	label.Y = c.Y
	label.Draw(c.PDF)

	value.Text = o.ID
	value.X = c.X + label.GetTextWidth(c.PDF)
	value.Y = c.Y
	value.Draw(c.PDF)
}
