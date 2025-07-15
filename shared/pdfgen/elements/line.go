package elements

import "github.com/phpdave11/gofpdf"

// Line represents a straight line with given width and color drawn between (x1,y1) and (x2,y2)
type Line struct {
	X1    float64
	Y1    float64
	X2    float64
	Y2    float64
	Width float64
	Color [3]int
}

func (l *Line) Draw(pdf *gofpdf.Fpdf) {
	pdf.SetLineWidth(l.Width)
	pdf.SetDrawColor(l.Color[0], l.Color[1], l.Color[2])
	pdf.Line(l.X1, l.Y1, l.X2, l.Y2)
}