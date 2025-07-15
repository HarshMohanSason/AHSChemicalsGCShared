package elements

import (
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/httpimg"
)

type Image struct {
	ImageURL string
	X        float64
	Y        float64
	Width    float64
	Height   float64
	Flow     bool
	Tpath    string
	Link     int
	LinkStr  string
}

func (i *Image) Draw(pdf *gofpdf.Fpdf) {
	httpimg.Register(pdf, i.ImageURL, "png")
	pdf.Image(i.ImageURL, i.X, i.Y, i.Width, i.Height, i.Flow, i.Tpath, i.Link, i.LinkStr)
}
