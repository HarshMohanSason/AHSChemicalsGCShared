package quickbooks

import (
	"regexp"
	"strconv"
	"strings"
)

//Available brands for the products
var brands = []string{"MicroTECH", "ProBlend"}

// ParseBrandName extracts the brand name from the given product name.
// It compares the input against a known list of brands in a case-insensitive way.
// Returns the matched brand or an empty string if no match is found.
func ParseBrandName(productName string) string {
	for _, brand := range brands {
		if strings.Contains(strings.ToLower(productName), strings.ToLower(brand)) {
			return brand
		}
	}
	return ""
}

// ParseProductName removes the brand name from the given product name, if present.
// This is useful for extracting the raw product title without branding.
// The removal is case-insensitive, and the result is trimmed of extra spaces.
func ParseProductName(productName string) string {
	brand := ParseBrandName(productName)
	if brand != "" {
		// Remove brand name from product name (case-insensitive)
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(brand))
		productName = re.ReplaceAllString(productName, "")
	}
	return strings.TrimSpace(productName)
}

// ParseSKU parses a SKU string formatted as "SKU-Size-Unit-PackOf".
// Example: "523423423-5-GAL-2" returns a map:
// {
//   "SKU": "523423423",
//   "Size": 5.0,
//   "SizeUnit": "GAL",
//   "PackOf": 2,
// }
// If parsing fails (due to wrong format), it returns default empty values.
func ParseSKU(productSKU string) map[string]any {
	splitString := strings.SplitN(productSKU, "-", 4)
	if len(splitString) == 4 {
		size, err := strconv.ParseFloat(splitString[1], 64)
		if err != nil {
			size = 0.0
		}
		packOf, err := strconv.Atoi(splitString[3])
		if err != nil {
			packOf = 0
		}
		parsedSKU := map[string]any{
			"SKU":      splitString[0],
			"Size":     size,
			"SizeUnit": splitString[2],
			"PackOf":   packOf,
		}
		return parsedSKU
	}
	// Return map with empty values if parsing fails
	return map[string]any{
		"SKU":      "",
		"Size":     0.0,
		"SizeUnit": "",
		"PackOf":   0,
	}
}
