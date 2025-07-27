package utils

import "time"

//HasDuplicates checks if a slice contains any duplicates.
//
//Parameters:
// - slice: The slice of type string to check
//
//Returns:
// - bool: True if the slice contains any duplicates, false otherwise
func HasDuplicates(slice []string) bool {
	seen := make(map[string]bool)
	for _, val := range slice {
		if seen[val] {
			return true // Duplicate found
		}
		seen[val] = true
	}
	return false
}

//FormatDate formats a time.Time object into a string in the format "January 2, 2006 at 3:04 PM"
func FormatDate(t time.Time) string {
	return t.Format("January 2, 2006 at 3:04 PM")
}