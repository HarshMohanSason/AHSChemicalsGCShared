package models

import (
	"errors"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/constants"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/utils"
)

// Map of valid roles.
var (
	roles = map[string]struct{}{
		constants.RoleSuperAdmin: {},
		constants.RoleAdmin:      {},
		constants.RoleUser:       {},
	}
)

// UserAccount represents a new user account created by the super-admin
type UserAccountCreate struct {
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Customers []string `json:"customers"`
	Brands    []string `json:"brands"`
	Role      string   `json:"role"`
}

// GetCustomClaimsMap returns the custom claims map for the user account with the role set to true.
func (c *UserAccountCreate) GetCustomClaimsMap() map[string]any {
	return map[string]any{
		c.Role: true,
	}
}

func (c *UserAccountCreate) ToFirestoreMap() map[string]any {
	return map[string]any{
		"name":      c.Name,
		"email":     c.Email,
		"customers": c.Customers,
		"brands":    c.Brands,
		"role":      c.Role,
	}
}

// Validate validates the basic user account details. Not doing a super strict validation here because
// firebase auth does it.
func (c *UserAccountCreate) Validate() error {
	if c.Name == "" {
		return errors.New("Name of the user cannot be empty")
	}
	if c.Email == "" {
		return errors.New("Email of the user cannot be empty")
	}
	if c.Password == "" {
		return errors.New("Password of the user cannot be empty")
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
	if c.Role == "" {
		return errors.New("Role of the user cannot be empty")
	}
	if _, ok := roles[c.Role]; !ok {
		return errors.New("Role of the user is not valid")
	}
	return nil
}

// Used when storing/retrieving from Firestore (no password)
type UserAccount struct {
	Name      string   `json:"name" firestore:"name"`
	Email     string   `json:"email" firestore:"email"`
	Customers []string `json:"customers" firestore:"customers"`
	Brands    []string `json:"brands" firestore:"brands"`
	Role      string   `json:"role" firestore:"role"`
}
