package validation

import (
	"errors"
	"regexp"
	"strings"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

func ValidateOrder(o *models.Order) error {
	if o.Uid == "" {
		return errors.New("No user id found when order was placed")
	}
	if len(o.Items) == 0 {
		return errors.New("No items found in order")
	}
	if o.TaxRate == 0 {
		return errors.New("No tax rate found in order")
	}
	if o.Customer.ID == "" {
		return errors.New("No customer found for this order")
	}
	if o.Customer.Name == "" {
		return errors.New("No customer name found for this order")
	}
	if o.TimeZone == "" {
		return errors.New("No timezone found for this order")
	}
	return validateSpecialInstructions(o.SpecialInstructions)
}

// validateSpecialInstructions checks that the special instructions are safe plain text without any HTML or script-like content.
func validateSpecialInstructions(specialInstructions string) error {
	if specialInstructions == "" {
		return nil
	}
	// Removing excess whitespace
	instructions := strings.TrimSpace(specialInstructions)

	// Reject if there are angle brackets (which might indicate HTML tags) or other unusual symbols
	illegalPattern := regexp.MustCompile(`[<>[\]{}^*~|\\]`)
	if illegalPattern.MatchString(instructions) {
		return errors.New("special instructions contain invalid characters")
	}

	if len(instructions) > 200 {
		return errors.New("special instructions are too long")
	}
	return nil
}
