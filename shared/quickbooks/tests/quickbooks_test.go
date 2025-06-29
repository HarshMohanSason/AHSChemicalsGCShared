package tests

import (
	"reflect"
	"testing"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
)

func TestParseBrandName(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Valid brand name",
			value:    "MICROTECHasdfsdaasdfaASSADASD 312@#!@asdffsfas2sadfas",
			expected: "MicroTECH",
		},
		{
			name:     "Empty brand name",
			value:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			brand := quickbooks.ParseBrandName(tc.value)
			if brand != tc.expected {
				t.Errorf("Error occurred, expected value: %s, got: %s", tc.expected, brand)
			}
		})
	}
}

func TestParseProductName(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Valid product name",
			value:    "MICROTECHasdfsdaasdfaASSADASD 312@#!@asdffsfas2sadfas",
			expected: "asdfsdaasdfaASSADASD 312@#!@asdffsfas2sadfas",
		},
		{
			name:     "Empty product name",
			value:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			productName := quickbooks.ParseProductName(tc.value)
			if productName != tc.expected {
				t.Errorf("Error occurred, expected value: %s, got: %s", tc.expected, productName)
			}
		})
	}
}

func TestParseSKU(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected map[string]any
	}{
		{
			name:  "Valid SKU",
			value: "523423423-5-GAL-2",
			expected: map[string]any{
				"SKU":      "523423423",
				"Size":     5.0,
				"SizeUnit": "GAL",
				"PackOf":   2,
			},
		},
		{
			name:  "Invalid input",
			value: "",
			expected: map[string]any{
				"SKU":      "",
				"Size":     0.0,
				"SizeUnit": "",
				"PackOf":   0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedSKU := quickbooks.ParseSKU(tc.value)
			if !reflect.DeepEqual(parsedSKU, tc.expected) {
				t.Errorf("Error occurred, expected value: %s, got: %s", tc.expected, parsedSKU)
			}
		})
	}
}
