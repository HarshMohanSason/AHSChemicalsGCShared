package elements

import "github.com/phpdave11/gofpdf"

//Rectangle drawn in the pdf. Can have border and fill color
type Rectangle struct {
	X           float64
	Y           float64
	Width       float64
	Height      float64
	Style       string
	BorderColor [3]int
	FillColor   [3]int
	LineWidth   float64
}

func (r *Rectangle) Draw(pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(r.BorderColor[0], r.BorderColor[1], r.BorderColor[2])
	pdf.SetLineWidth(r.LineWidth)
	pdf.SetFillColor(r.FillColor[0], r.FillColor[1], r.FillColor[2])
	pdf.Rect(r.X, r.Y, r.Width, r.Height, r.Style)
}