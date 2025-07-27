package firestore_testutils

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

func CreateMockCustomer() *models.Customer {
	return &models.Customer{
		ID:   "12345",
		QBID: "123",
		Name: "Test Customer Name 1",
		Email: "g7vP7@example.com",
		Phone: "1234567890",
		Address1: "2040 N preisker lane", 
		City: "Las Vegas",
		State: "NV",
		Zip: "12345",
		Country: "USA",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}