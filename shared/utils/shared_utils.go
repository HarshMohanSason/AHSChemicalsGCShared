package utils

import (
	"bytes"
	"image"
	"io"
	"math"
	"mime/multipart"
	"os"
	"strings"
)

//HasDuplicateStrings checks if a slice of string contains any duplicates.
//
//Parameters:
// - slice: The slice of type string to check
//
//Returns:
// - bool: True if the slice contains any duplicates, false otherwise
//
func HasDuplicateStrings(slice []string) bool {
	seen := make(map[string]bool)
	for _, val := range slice {
		formattedVal := strings.ToLower(strings.TrimSpace(val))
		formattedVal = strings.ReplaceAll(formattedVal, " ", "")
		if seen[formattedVal] {
			return true // Duplicate found
		}
		seen[formattedVal] = true
	}
	return false
}

//AreEqualStringSlices checks if two slices of strings are equal.
//
//Parameters:
// - a: The first slice of type string
// - b: The second slice of type string
//
//Returns:
// - bool: True if the slices are equal, false otherwise
//
func AreEqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		v = strings.ToLower(strings.TrimSpace(v))
		v = strings.ReplaceAll(v, " ", "")
		b[i] = strings.ToLower(strings.TrimSpace(b[i]))
		b[i] = strings.ReplaceAll(b[i], " ", "")
		if v != b[i] {
			return false
		}
	}
	return true
}

//DetectImageType detects the type of an image. 
//
//Parameters:
// - imageBytes: The bytes of the image
//
//Returns:
// - string: The type of the image
func DetectImageType(imageBytes []byte) string {
	_, format, err := image.DecodeConfig(bytes.NewReader(imageBytes))
	if err != nil {
		return "png" 
	}
	return format
}

//GetImageBytesFromMultiPart gets the bytes of an image from a multipart file.
//
//Parameters:
// - file: The multipart file
//
//Returns:
// - []byte: The bytes of the image
// - error: The error if any
//
func GetImageBytesFromMultiPart(file multipart.File) ([]byte, error) {
	defer file.Close()
	return io.ReadAll(file)
}

//CreateMultipartFile creates a multipart file from a file path.
//
//Parameters:
// - path: The path of the file
//
//Returns: 
// - multipart.File: The multipart file
// - error: The error if any
//
func CreateMultipartFile(path string) (multipart.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

//Rounds to a specified number of decimal places.
func RoundToDecimals(val float64, place float64) float64 {
	return math.Round(val * 10*(place)) / (10*place)
}