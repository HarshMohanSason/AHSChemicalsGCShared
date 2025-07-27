package elements

import (
	"bytes"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/httpimg"
)

type Image struct {
	ImageURL   string
	ImageBytes []byte
	X          float64
	Y          float64
	Width      float64
	Height     float64
	Flow       bool
	Tpath      string
	Link       int
	LinkStr    string
}

func (i *Image) detectImageType() string {
	_, format, err := image.DecodeConfig(bytes.NewReader(i.ImageBytes))
	if err != nil {
		return "png" //default
	}
	return format
}

// Draws the image from []bytes
func (i *Image) Draw(pdf *gofpdf.Fpdf) {
	options := gofpdf.ImageOptions{
		ImageType:             strings.ToUpper(i.detectImageType()),
		ReadDpi:               true,
		AllowNegativePosition: false,
	}
	imageName := "img" + fmt.Sprintf("%d", time.Now().UnixNano()) 
	pdf.RegisterImageOptionsReader(imageName, options, bytes.NewReader(i.ImageBytes))
	
	// Draw the image on the canvas
	pdf.ImageOptions(
		imageName,
		i.X, i.Y, i.Width, i.Height, i.Flow, options, i.Link, i.LinkStr)
}

func (i *Image) DrawFromURL(pdf *gofpdf.Fpdf) {
	httpimg.Register(pdf, i.ImageURL, "png")
	pdf.Image(i.ImageURL, i.X, i.Y, i.Width, i.Height, i.Flow, i.Tpath, i.Link, i.LinkStr)
}
