package firestore_testutils

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/constants"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

func CreateMockOrder() *models.Order {
	now := time.Now()

	return &models.Order{
		ID:                  "mock-order-id-001",
		Customer:            *CreateMockCustomer(),
		Uid:                 "mock-user-uid-123",
		SpecialInstructions: "Leave at the front desk.",
		Items:               CreateMultipleMockProducts(20),
		MinimalItems:        []models.ItemMinimal{},
		TaxRate:             0.08, 
		SubTotal:            500.00,
		TaxAmount:           500.00 * 0.08,
		Total:               500.00 + (500.00 * 0.08),
		Status:              constants.OrderStatusPending,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}