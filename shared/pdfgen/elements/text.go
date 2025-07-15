package elements

import "github.com/phpdave11/gofpdf"

type Text struct {
	Text  string
	Font  string
	Style string
	X     float64
	Y     float64
	Size  float64
	Color [3]int
}

func (t *Text) ChangeCurrentTextStyle(font string, style string, size float64) {
	t.Font = font
	t.Style = style
	t.Size = size
}

func (t *Text) ApplyTextStyle(pdf *gofpdf.Fpdf) {
	pdf.SetFont(t.Font, t.Style, t.Size)
	pdf.SetTextColor(t.Color[0], t.Color[1], t.Color[2])
}

func (t *Text) GetTextWidth(pdf *gofpdf.Fpdf) float64 {
	t.ApplyTextStyle(pdf)
	return pdf.GetStringWidth(t.Text)
}

func (t *Text) GetTextHeight(pdf *gofpdf.Fpdf) float64 {
	t.ApplyTextStyle(pdf)
	fontSize, _ := pdf.GetFontSize()
	return pdf.PointConvert(fontSize)
}

func (t *Text) GetMultiLineHeight(pdf *gofpdf.Fpdf, allowedWidth float64, lineHeight float64) float64 {
	t.ApplyTextStyle(pdf)
	lines := pdf.SplitLines([]byte(t.Text), allowedWidth)
	return float64(len(lines)) * lineHeight
}

func (t *Text) Draw(pdf *gofpdf.Fpdf) {
	t.ApplyTextStyle(pdf)
	pdf.Text(t.X, t.Y, t.Text)
}

func (t *Text) DrawMultipleLines(pdf *gofpdf.Fpdf, allowedWidth float64, align string) float64 {
	t.ApplyTextStyle(pdf)
	lineHeight := t.GetTextHeight(pdf)
	lines := pdf.SplitLines([]byte(t.Text), allowedWidth)

	for i, line := range lines {
		text := string(line)
		var x float64

		switch align {
		case "center":
			lineWidth := pdf.GetStringWidth(text)
			x = t.X + (allowedWidth/2 - lineWidth/2)
		case "right":
			lineWidth := pdf.GetStringWidth(text)
			x = t.X + (allowedWidth - lineWidth)
		default: // "left"
			x = t.X
		}
		pdf.Text(x, t.Y+float64(i)*lineHeight, text)
	}
	return float64(len(lines)) * lineHeight
}