package utils

import (
	"errors"
	"regexp"
)

func ValidateContactUsForm(form map[string]any) error {
	err := ValidateContactUsFormEmail(form["email"].(string))
	if err != nil {
		return err
	}
	err = ValidateContactUsFormName(form["name"].(string))
	if err != nil {
		return err
	}
	err = ValidateContactUsFormPhone(form["phone"].(string))
	if err != nil {
		return err
	}
	err = ValidateContactUsFormLocation(form["location"].(string))
	if err != nil {
		return err
	}
	err = ValidateContactUsFormMessage(form["message"].(string))
	if err != nil {
		return err
	}
	return nil
}

func ValidateContactUsFormEmail(email string) error {
	if email == "" {
		return errors.New("Email is required")
	}
	reg, err := regexp.Compile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	if err != nil {
		return err
	}
	if !reg.MatchString(email) {
		return errors.New("Invalid email address entered. Please enter the email in the correct format")
	}
	return nil
}

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

func ValidateContactUsFormPhone(phone string) error {
	if phone == "" {
		return errors.New("Phone number is required")
	}
	reg, err := regexp.Compile(`^[0-9]{10}$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(phone) {
		return errors.New("Invalid phone number entered.Only need the 10 digits of the number")
	}
	return nil
}

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

func ValidateContactUsFormMessage(message string) error {
	if message == "" {
		return errors.New("Message is required")
	}
	reg, err := regexp.Compile(`^[a-zA-Z0-9\s.,'"!?;:\-_()\n\r]{10,1000}$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(message) {
		return errors.New("Invalid message entered. Message should be between 10 and 1000 characters with no symbols or special characters")
	}
	return nil
}