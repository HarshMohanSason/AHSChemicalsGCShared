package purchase_order

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/elements"
)

func DrawVendorShipToSectionHeading(text string, c *canvas.Canvas){
	const height = 5
	const width = 70

	rectElement := elements.Rectangle{
		X: c.X,
		Y: c.Y,
		Width: width,
		Height: height,
		Style: "F",
		FillColor: canvas.PrimaryBlue,
	} 
	rectElement.Draw(c.PDF)

	textElement := elements.Text{
		Text: text,
		Font: "Arial",
		Style: "B",
		X: c.X + 5,
		Y: c.Y,
		Size: 10,
		Color: canvas.White,
	}
	textElement.ApplyTextStyle(c.PDF)
	textHeight := textElement.GetTextHeight(c.PDF)
	textElement.Y += (height + textHeight)/2 - 0.5
	textElement.Draw(c.PDF)
}