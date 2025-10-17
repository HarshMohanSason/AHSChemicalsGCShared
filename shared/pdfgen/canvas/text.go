package canvas

import (
	"github.com/phpdave11/gofpdf"
)

type Text struct {
	Content string
	Font    string
	Style   string
	X       float64
	Y       float64
	Size    float64
	Color   [3]int
}

func (t *Text) SetStyle(style string)     { t.Style = style }
func (t *Text) SetFont(font string)       { t.Font = font }
func (t *Text) SetSize(size float64)      { t.Size = size }
func (t *Text) SetContent(content string) { t.Content = content }
func (t *Text) SetX(x float64)            { t.X = x }
func (t *Text) SetY(y float64)            { t.Y = y }
func (t *Text) SetColor(color [3]int) {
	t.Color[0] = color[0]
	t.Color[1] = color[1]
	t.Color[2] = color[2]
}

func (t *Text) ApplyTextStyle(pdf *gofpdf.Fpdf) {
	pdf.SetFont(t.Font, t.Style, t.Size)
	pdf.SetTextColor(t.Color[0], t.Color[1], t.Color[2])
}

// Only used for single line text
func (t *Text) GetTextHeight(pdf *gofpdf.Fpdf) float64 {
	t.ApplyTextStyle(pdf)
	_, unitSize := pdf.GetFontSize()
	return unitSize
}

func (t *Text) GetWidth(pdf *gofpdf.Fpdf) float64 {
	t.ApplyTextStyle(pdf)
	return pdf.GetStringWidth(t.Content)
}

func (t *Text) GetDescent(pdf *gofpdf.Fpdf) float64 {
	t.ApplyTextStyle(pdf)
	return float64(pdf.GetFontDesc(t.Font, t.Style).Descent)
}

func (t *Text) GetAscent(pdf *gofpdf.Fpdf) float64 {
	t.ApplyTextStyle(pdf)
	return float64(pdf.GetFontDesc(t.Font, t.Style).Ascent)
}

func (t *Text) GetMultiTextHeight(pdf *gofpdf.Fpdf, drawingWidth float64) float64 {
	t.ApplyTextStyle(pdf)
	lineHeight := t.GetTextHeight(pdf)
	lines := pdf.SplitLines([]byte(t.Content), drawingWidth)
	return float64(len(lines)) * lineHeight
}

/* Canvas methods */

func (c *Canvas) DrawSingleLineText(text *Text) {
	text.ApplyTextStyle(c.PDF)
	c.PDF.Text(text.X, text.Y, text.Content)
}

// Draws multiple lines of text using the same font, style and size.
// Does not support checking new page after each line because this function is used
// to draw multiple lines of text in a single row. The caller is responsible for
// checking new page when creating the row by first checking the row height.
func (c *Canvas) DrawMultipleLines(t *Text, allowedWidth float64, align string) {
	t.ApplyTextStyle(c.PDF)
	lineHeight := t.GetTextHeight(c.PDF)
	lines := c.PDF.SplitLines([]byte(t.Content), allowedWidth)

	for i, line := range lines {
		text := string(line)
		var x float64
		switch align {
		case "C":
			lineWidth := c.PDF.GetStringWidth(text)
			x = t.X + (allowedWidth/2 - lineWidth/2)
		case "L":
			lineWidth := c.PDF.GetStringWidth(text)
			x = t.X + (allowedWidth - lineWidth)
		default: // "left"
			x = t.X
		}
		c.PDF.Text(x, t.Y+float64(i)*lineHeight, text)
	}
}

func (c *Canvas) DrawLabelWithSingleLineText(t *Text, value string) {
	t.ApplyTextStyle(c.PDF)
	c.DrawSingleLineText(t)

	labelWidth := c.PDF.GetStringWidth(t.Content)

	t.SetContent(value)
	t.SetStyle("")
	t.SetX(t.X + labelWidth + 2) // 2 is the padding
	c.DrawSingleLineText(t)
}

func (c *Canvas) DrawTextInColoredRect(t *Text, rect *Rectangle, align string) {
	c.DrawRectangle(rect)
	switch align {
	case "center":
		t.SetX(rect.X + (rect.Width/2 - c.PDF.GetStringWidth(t.Content)/2))
	case "right":
		t.SetX(rect.X + rect.Width - c.PDF.GetStringWidth(t.Content) - 5) // 5 is the padding
	default: //"left"
		t.SetX(rect.X + 5) // 5 is the padding
	}
	baselineAdjustment := t.GetTextHeight(c.PDF) * 0.3 // y is the baseline of the text not the center. Adjusting it by 0.3
	t.SetY(rect.Y + rect.Height/2 + baselineAdjustment)
	c.DrawSingleLineText(t)
}
