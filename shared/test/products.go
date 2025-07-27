package firestore_testutils

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

func CreateMultipleMockProducts(length int)[]models.Product{
	var products []models.Product
	for range length {
		products = append(products, CreateMockProduct())
	}
	return products
}

func CreateMockProduct() models.Product {
	return models.Product{
		ID:        "mock-id-123",
		QBID:      "qb-001",
		Brand:     "TestBrand",
		Name:      "Test Product",
		SKU:       "TESTSKU123",
		Size:      1.0,
		SizeUnit:  "GAL",
		PackOf:    1,
		Hazardous: false,
		Category:  "Chemicals",
		Price:     100.0,
		Desc:      "This is a mock test product for development use.",
		Slug:      "test-product",
		NameKey:   "testproduct",
		Quantity:  100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}