package models

import (
	"fmt"
	"time"
)

// Customer represents a customer refined and cleaner version of the QBCustomer struct. 
type Customer struct {
	ID        string    `json:"id" firestore:"omitempty"`
	QBID      string    `json:"qbid" firestore:"qbid"` //quickbooks ID
	Name      string    `json:"name" firestore:"name"`
	Email     string    `json:"email" firestore:"email"`
	Phone     string    `json:"phone" firestore:"phone"`
	Address1  string    `json:"address1" firestore:"address1"`
	City      string    `json:"city" firestore:"city"`
	State     string    `json:"state" firestore:"state"`
	Zip       string    `json:"zip" firestore:"zip"`
	Country   string    `json:"country" firestore:"country"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// FormatAddress2 returns City, State, Zip
func (c *Customer) FormatAddress2() string {
	return fmt.Sprintf("%s, %s %s", c.City, c.State, c.Zip)
}