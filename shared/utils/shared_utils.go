package utils

import (
	"strings"
	"time"
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

//FormatDate formats a time.Time object into a string in the format "January 2, 2006 at 3:04 PM"
func FormatDate(t time.Time) string {
	return t.Format("January 2, 2006 at 3:04 PM")
}