package utils

import (
	"errors"
	"regexp"
)

// ValidateContactUsForm validates all fields of the contact us form.
// The `form` map is expected to contain the following keys:
// "email", "name", "phone", "location", and "message", each with string values.
//
// Returns:
//   - error: If any of the individual validations fail, the corresponding error is returned.
func ValidateContactUsForm(form map[string]any) error {
	if err := ValidateContactUsFormEmail(form["email"].(string)); err != nil {
		return err
	}
	if err := ValidateContactUsFormName(form["name"].(string)); err != nil {
		return err
	}
	if err := ValidateContactUsFormPhone(form["phone"].(string)); err != nil {
		return err
	}
	if err := ValidateContactUsFormLocation(form["location"].(string)); err != nil {
		return err
	}
	if err := ValidateContactUsFormMessage(form["message"].(string)); err != nil {
		return err
	}
	return nil
}

// ValidateContactUsFormEmail checks that the email is non-empty and conforms to a valid email format.
//
// Parameters:
//   - email: The email string to validate.
//
// Returns:
//   - error: If the email is missing or improperly formatted.
func ValidateContactUsFormEmail(email string) error {
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

// ValidateContactUsFormName checks that the name is non-empty and consists of two words (first and last name).
//
// Parameters:
//   - name: The full name string to validate.
//
// Returns:
//   - error: If the name is missing or doesn't follow the format "Firstname Lastname".
func ValidateContactUsFormName(name string) error {
	if name == "" {
		return errors.New("Name is required")
	}
	reg, err := regexp.Compile(`^[a-zA-Z]+\s[a-zA-Z]+$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(name) {
		return errors.New("Invalid name entered. Please enter the name in the correct format")
	}
	return nil
}

// ValidateContactUsFormPhone checks that the phone number is exactly 10 digits.
//
// Parameters:
//   - phone: The phone number string to validate.
//
// Returns:
//   - error: If the phone is missing or doesn't contain exactly 10 digits.
func ValidateContactUsFormPhone(phone string) error {
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

// ValidateContactUsFormLocation checks that the location is non-empty and between 10 to 50 characters long,
// allowing alphanumeric characters, spaces, commas, periods, underscores, and hyphens.
//
// Parameters:
//   - location: The location string to validate.
//
// Returns:
//   - error: If the location is missing or contains disallowed characters.
func ValidateContactUsFormLocation(location string) error {
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

// ValidateContactUsFormMessage checks that the message is non-empty and between 10 to 300 characters,
// allowing letters, numbers, spaces, and selected punctuation characters.
//
// Parameters:
//   - message: The message string to validate.
//
// Returns:
//   - error: If the message is missing or contains disallowed characters.
func ValidateContactUsFormMessage(message string) error {
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