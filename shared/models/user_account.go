package models

import (
	"errors"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/utils"
)

// UserAccount represents a new user account created by the admin
type UserAccount struct {
	UID         string      `json:"uid"`
	Name        string      `json:"name"`
	PhoneNumber PhoneNumber `json:"phoneNumber"`
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	Customers   []string    `json:"customers" firestore:"customers"`
	Brands      []string    `json:"brands" firestore:"brands"`
}

// Validate validates the user account. Email and password are validated by firebase authentication, so only the necessary fields are validated
func (c *UserAccount) Validate() error {
	if c.Name == "" {
		return errors.New("Name of the user cannot be empty")
	}
	if err := c.PhoneNumber.Validate(); err != nil {
		return err
	}
	if len(c.Customers) == 0 {
		return errors.New("At least one customer is required for the user")
	}
	if utils.HasDuplicateStrings(c.Customers) {
		return errors.New("Customers of the user cannot have duplicates")
	}
	if len(c.Brands) == 0 {
		return errors.New("Brands of the user cannot be empty")
	}
	if utils.HasDuplicateStrings(c.Brands) {
		return errors.New("Brands of the user cannot have duplicates")
	}
	return nil
}

func (c *UserAccount) ToUserAccountUpdate() *UserAccountUpdate {
	return &UserAccountUpdate{
		UID:       c.UID,
		Brands:    c.Brands,
		Customers: c.Customers,
	}
}

// UpdatedAccount represents an updated account. It represents the firestore doucument for the user more or less
type UserAccountUpdate struct {
	UID       string   `json:"uid"`
	Brands    []string `json:"brands"`
	Customers []string `json:"customers"`
}

func (u *UserAccountUpdate) Validate() error {
	if u.UID == "" {
		return errors.New("No uid found for the updated account")
	}
	if len(u.Brands) == 0 {
		return errors.New("No brands found for the updated account")
	}
	if utils.HasDuplicateStrings(u.Brands) {
		return errors.New("Duplicate brands found for the updated account")
	}
	if len(u.Customers) == 0 {
		return errors.New("No customers found for the updated account")
	}
	if utils.HasDuplicateStrings(u.Customers) {
		return errors.New("Duplicate customers found for the updated account")
	}
	return nil
}

func (u *UserAccountUpdate) ToMap() map[string]any {
	return map[string]any{
		"customers": u.Customers,
		"brands":    u.Brands,
	}
}
