package mocks

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

func CreateMockProduct() *models.Product {
	return &models.Product{
		ID:            "prod-001",
		IsActive:      true,
		Brand:         "Acme",
		Name:          "Industrial Cleaner",
		SKU:           "ACM-CLN-500ML",
		Size:          500.0,
		SizeUnit:      "ml",
		PackOf:        12,
		Hazardous:     true,
		Category:      "Chemicals",
		Price:         129.99,
		PurchasePrice: 50.00,
		Desc:          "Highly effective industrial cleaning solution.",
		Slug:          "industrial-cleaner-500ml",
		NameKey:       "industrialcleaner",
		Quantity:      200,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func CreateMockProducts(val int) []*models.Product {
	products := make([]*models.Product, val)
	for i := range val {
		products[i] = CreateMockProduct()
	}
	return products
}
