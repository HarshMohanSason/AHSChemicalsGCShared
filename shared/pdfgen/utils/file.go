package utils

import (
	"bytes"

	"github.com/phpdave11/gofpdf"
)

//Returns the PDF output in bytes
func GetGeneratedPDF(pdf *gofpdf.Fpdf) ([]byte, error) {
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}