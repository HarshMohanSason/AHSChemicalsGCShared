package purchase_order

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/canvas"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/pdfgen/elements"
)

func DrawCommentsorSpecialInstructions(specialInstructions string, c *canvas.Canvas) {
	
	//Label
	text := elements.Text{
		Text:  "Comments or Special Instructions",
		Font:  "Arial",
		Style: "B",
		X:     c.X,
		Y:     c.Y,
		Size:  10,
		Color: canvas.Black,
	}	
	text.Draw(c.PDF)
	textHeight := text.GetTextHeight(c.PDF)
	c.IncY(textHeight + 1)

	text.Y = c.Y; text.Style = "";
	text.Text = specialInstructions
	text.DrawMultipleLines(c.PDF, 65, "left")
}