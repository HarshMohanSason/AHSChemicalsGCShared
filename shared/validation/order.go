package validation

import (
	"errors"
	"regexp"
	"strings"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

// ValidateSpecialInstructions checks that the special instructions are safe plain text without any HTML or script-like content.
func ValidateSpecialInstructions(o *models.Order) error {
	if o.SpecialInstructions == "" {
		return nil
	}
	// Removing excess whitespace
	instructions := strings.TrimSpace(o.SpecialInstructions)

	// Reject if there are angle brackets (which might indicate HTML tags) or other unusual symbols
	illegalPattern := regexp.MustCompile(`[<>[\]{}$%^*~|\\]`)
	if illegalPattern.MatchString(instructions) {
		return errors.New("special instructions contain invalid characters")
	}

	if len(instructions) > 200 {
		return errors.New("special instructions are too long")
	}
	return nil
}


// CheckIfOrderItemsChanged checks if the items in the order have any changes.
//
// Parameters:
//   - editedOrder: Pointer to an Order object containing the edited order.
//   - originalOrder: Pointer to an Order object containing the original order.
//
// Returns:
//   - bool: A boolean value indicating whether the items in the order have any changes.
func CheckIfOrderItemsChanged(editedOrder, originalOrder *models.Order) bool {

	for i, item := range editedOrder.Items {
		if item.Price != originalOrder.Items[i].Price {
			return true
		}
		if item.Quantity != originalOrder.Items[i].Quantity {
			return true
		}
	}
	return false
}