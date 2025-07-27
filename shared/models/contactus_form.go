package models

import (
	"time"
)

// ContactUsForm represents the structure of a contact submission form.
type ContactUsForm struct {
	Email     string    `json:"email" firestore:"email"`    
	Name      string    `json:"name" firestore:"name"`     
	Phone     string    `json:"phone" firestore:"phone"`    
	Location  string    `json:"location" firestore:"location"` 
	Message   string    `json:"message" firestore:"message"`  
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}