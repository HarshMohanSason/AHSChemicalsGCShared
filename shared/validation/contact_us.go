// Validation package handles complex validations in the shared package
package validation

import (
	"errors"
	"regexp"
	"strings"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

// ValidateContactUsForm performs validation on each field of the ContactUsForm.
// It returns the first encountered error to signal a failed validation.
//
// Validation rules:
// - Email must be a valid email address format.
// - Name must contain only alphabetic characters and spaces.
// - Phone must contain only digits and be between 7 to 15 characters.
// - Location must consist of alphanumeric characters and basic punctuation.
// - Message must be plain text with no special characters or HTML tags.
func ValidateContactUsForm(c *models.ContactUsForm) error {
	if err := validateEmail(c.Email); err != nil {
		return err
	}
	if err := validateName(c.Name); err != nil {
		return err
	}
	if err := validatePhone(c.Phone); err != nil {
		return err
	}
	if err := validateLocation(c.Location); err != nil {
		return err
	}
	if err := validateMessage(c.Message); err != nil {
		return err
	}
	return nil
}

// validateEmail checks that the email is non-empty and conforms to a valid email format.
//
// Parameters:
//   - email: The email string to validate.
//
// Returns:
//   - error: If the email is missing or improperly formatted.
func validateEmail(email string) error {
	if email == "" {
		return errors.New("Email is required")
	}
	reg, err := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(email) {
		return errors.New("Invalid email address entered. Please enter the email in the correct format")
	}
	return nil
}

// validateName checks that the name is non-empty and consists of two words (first and last name).
//
// Parameters:
//   - name: The full name string to validate.
//
// Returns:
//   - error: If the name is missing or doesn't follow the format "Firstname Lastname".
func validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("Name is required")
	}
	reg, err := regexp.Compile(`^[A-Za-z]+(?:\s+[A-Za-z]+)+$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(name) {
		return errors.New("Invalid name entered. Please enter your full name (first and last)")
	}
	return nil
}

// validatePhone checks that the phone number is exactly 10 digits.
//
// Parameters:
//   - phone: The phone number string to validate.
//
// Returns:
//   - error: If the phone is missing or doesn't contain exactly 10 digits.
func validatePhone(phone string) error {
	if phone == "" {
		return errors.New("Phone number is required")
	}
	reg, err := regexp.Compile(`^[0-9]{10}$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(phone) {
		return errors.New("Invalid phone number entered. Only need the 10 digits of the number")
	}
	return nil
}

// validateLocation checks that the location is non-empty and between 10 to 50 characters long,
// allowing alphanumeric characters, spaces, commas, periods, underscores, and hyphens.
//
// Parameters:
//   - location: The location string to validate.
//
// Returns:
//   - error: If the location is missing or contains disallowed characters.
func validateLocation(location string) error {
	if location == "" {
		return errors.New("Location is required")
	}
	reg, err := regexp.Compile(`^[a-zA-Z0-9 ,._-]{10,50}$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(location) {
		return errors.New("Invalid location entered. Location can only be between 10 and 50 characters with no symbols or special characters")
	}
	return nil
}

// validateMessage checks that the message is non-empty and between 10 to 300 characters,
// allowing letters, numbers, spaces, and selected punctuation characters.
//
// Parameters:
//   - message: The message string to validate.
//
// Returns:
//   - error: If the message is missing or contains disallowed characters.
func validateMessage(message string) error {
	if message == "" {
		return errors.New("Message is required")
	}
	reg, err := regexp.Compile(`^[a-zA-Z0-9\s.,'"!?;:\-_()\n\r]{10,300}$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(message) {
		return errors.New("Invalid message entered. Message should be between 10 and 300 characters with no symbols or special characters")
	}
	return nil
}