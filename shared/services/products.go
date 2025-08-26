package services

import (

	"cloud.google.com/go/firestore"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

// GetUpdatedProductDetails is similar to the GetUpdatedCustomerDetails. Returns the updated product details.
// Used for creating a map of what has changed when the customer is updated in quickbooks and is detected
// by the webhook.
func GetUpdatedProductDetails(newProduct, oldProduct *models.Product) map[string]any {
	if newProduct == nil || oldProduct == nil {
		return nil
	}

	changedValues := make(map[string]any)

	if newProduct.IsActive != oldProduct.IsActive {
		changedValues["isActive"] = newProduct.IsActive
	}
	if newProduct.Brand != oldProduct.Brand {
		changedValues["brand"] = newProduct.Brand
	}
	if newProduct.Name != oldProduct.Name {
		changedValues["name"] = newProduct.Name
	}
	if newProduct.SKU != oldProduct.SKU {
		changedValues["sku"] = newProduct.SKU
	}
	if newProduct.Size != oldProduct.Size {
		changedValues["size"] = newProduct.Size
	}
	if newProduct.SizeUnit != oldProduct.SizeUnit {
		changedValues["sizeUnit"] = newProduct.SizeUnit
	}
	if newProduct.PackOf != oldProduct.PackOf {
		changedValues["packOf"] = newProduct.PackOf
	}
	if newProduct.Hazardous != oldProduct.Hazardous {
		changedValues["hazardous"] = newProduct.Hazardous
	}
	if newProduct.Category != oldProduct.Category {
		changedValues["category"] = newProduct.Category
	}
	if newProduct.Price != oldProduct.Price {
		changedValues["price"] = newProduct.Price
	}
	if newProduct.Desc != oldProduct.Desc {
		changedValues["desc"] = newProduct.Desc
	}
	if newProduct.Slug != oldProduct.Slug {
		changedValues["slug"] = newProduct.Slug
	}
	if newProduct.NameKey != oldProduct.NameKey {
		changedValues["nameKey"] = newProduct.NameKey
	}
	if newProduct.Quantity != oldProduct.Quantity {
		changedValues["quantity"] = newProduct.Quantity
	}

	if len(changedValues) == 0 {
		return nil
	}
	changedValues["updatedAt"] = firestore.ServerTimestamp

	return changedValues
}