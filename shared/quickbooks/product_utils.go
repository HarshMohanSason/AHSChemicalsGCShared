package quickbooks

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//Available brands for the products
var brands = []string{"MicroTECH", "ProBlend"}

// ParseBrandName extracts the brand name from a product name formatted as:
// "BrandName - Product Name".
// It matches against a predefined list of brands in a case-insensitive manner.
// Returns the matched brand name, or an empty string if no known brand is found.
func ParseBrandName(productName string) string {
	for _, brand := range brands {
		if strings.Contains(strings.ToLower(productName), strings.ToLower(brand)) {
			return brand
		}
	}
	return ""
}

// ParseProductName removes the brand name from the given product name,
// assuming the format is: "BrandName - Product Name".
// This is necessary because QuickBooks does not allow brand names,
// so they are embedded in the product name string.
// The function performs a case-insensitive match and removal,
// and trims any leading/trailing spaces from the result.
func ParseProductName(productName string) string {
	brand := ParseBrandName(productName)
	if brand != "" {
		// Remove brand name from product name (case-insensitive)
		re, err := regexp.Compile(`(?i)` + regexp.QuoteMeta(brand))
		if err != nil{
			return productName
		}
		productName = re.ReplaceAllString(productName, "")
	}
	return strings.TrimSpace(productName)
}

// ParseSKU extracts structured product details from a SKU string formatted as:
// "SKU-Size-Unit-PackOf".
// This format is used because platforms like QuickBooks do not support
// separate fields for size, units, or packaging information, so they are
// embedded directly in the SKU string in quickbooks.
//
// Example input:
//   "523423423-5-GAL-2"
//
// Returns a map:
//   {
//     "SKU":      "523423423",
//     "Size":     5.0,
//     "SizeUnit": "GAL",
//     "PackOf":   2,
//   }
//
// If parsing fails due to unexpected format, default zero or empty values are returned.
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

// CreateProductSlugWithNameKey generates a URL-friendly slug and a normalized name key
// for a product, using its name and unique product ID.
//
// The slug is created by:
// - Converting the product name to lowercase
// - Replacing all non-alphanumeric characters with hyphens
// - Trimming leading/trailing hyphens
// - Appending the product ID to ensure uniqueness
//
// The function returns a map with:
// - "Slug": the full slug string (e.g., "super-cleaner-500ml-abc123")
// - "NameKey": the normalized product name key (e.g., "super-cleaner-500ml")
//   which can be used to group product variants.
//
// Parameters:
// - productName: the name of the product (e.g., "Super Cleaner 500ml")
// - productID: the unique ID of the product (e.g., "abc123")
//
// Returns:
// - map[string]string with keys "Slug" and "NameKey"
func CreateProductSlugWithNameKey(productName string, productID string) map[string]string{
	
	slug := strings.ToLower(productName)
	re := regexp.MustCompile(`[^a-z0-9]+`) //Only lowercase letters and numbers in a slug

	slug = re.ReplaceAllString(slug, "-") //Replace everthing with a hyphen
	slug = strings.Trim(slug, "-") //Cleaning any leading or trailing spaces if any present
	nameKey := slug
	slug = fmt.Sprintf("%s-%s",slug, productID)
	
	productSlugWithNameKey := map[string]string{
		"Slug": slug, 
		"NameKey": nameKey,
	}
	return productSlugWithNameKey
}
