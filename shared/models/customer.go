// package models contains all the models used to represent data stored in firestore.
package models

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

// Customer represents a customer refined and cleaner version of the QBCustomer struct.
type Customer struct {
	ID        string    `json:"id" firestore:"id"`
	IsActive  bool      `json:"is_active" firestore:"isActive"`
	Name      string    `json:"name" firestore:"name"`
	Email     string    `json:"email" firestore:"email"`
	Phone     string    `json:"phone" firestore:"phone"`
	Address1  string    `json:"address1" firestore:"address1"`
	City      string    `json:"city" firestore:"city"`
	State     string    `json:"state" firestore:"state"`
	Zip       string    `json:"zip" firestore:"zip"`
	Country   string    `json:"country" firestore:"country"`
	CreatedAt time.Time `json:"created_at" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updatedAt"`
}

// FormatAddress2 returns City, State, Zip
func (c *Customer) FormatAddress2() string {
	return fmt.Sprintf("%s, %s %s", c.City, c.State, c.Zip)
}

func (c *Customer) ToMap() map[string]any {
	return map[string]any{
		"id":        c.ID,
		"isActive":  c.IsActive,
		"name":      c.Name,
		"email":     c.Email,
		"phone":     c.Phone,
		"address1":  c.Address1,
		"city":      c.City,
		"state":     c.State,
		"zip":       c.Zip,
		"country":   c.Country,
		"createdAt": firestore.ServerTimestamp,
		"updatedAt": firestore.ServerTimestamp,
	}
}
