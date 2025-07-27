package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/phpdave11/gofpdf"
)

//For debugging purposes
func GeneratePDFFileInPath(pdf *gofpdf.Fpdf, fileName string) error {
	err := os.MkdirAll("./shared/pdfgen/generated", os.ModePerm)
	if err != nil {
		return err
	}

	formattedPath := fmt.Sprintf("./shared/pdfgen/generated/%s.pdf", fileName)
	err = pdf.OutputFileAndClose(formattedPath)
	if err != nil {
		return err
	}

	return nil
}

func GeneratePDFBase64 (pdf *gofpdf.Fpdf) (string, error) {
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Str, nil
}

func MultipartFileToBytes(file multipart.File) ([]byte, error) {
	defer file.Close()
	return io.ReadAll(file)
}