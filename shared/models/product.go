package models

import (
	"fmt"
	"strings"
	"time"
)

type Product struct {
	ID            string    `json:"id" firestore:"id"`
	IsActive      bool      `json:"is_active" firestore:"isActive"`
	Brand         string    `json:"brand" firestore:"brand"`
	Name          string    `json:"name" firestore:"name"`
	SKU           string    `json:"sku" firestore:"sku"`
	Size          float64   `json:"size" firestore:"size"`
	SizeUnit      string    `json:"sizeUnit" firestore:"sizeUnit"`
	PackOf        int       `json:"packOf" firestore:"packOf"`
	Hazardous     bool      `json:"hazardous" firestore:"hazardous"`
	Category      string    `json:"category" firestore:"category"`
	Price         float64   `json:"price" firestore:"price"`
	PurchasePrice float64   `json:"purchasePrice" firestore:"purchasePrice"`
	Desc          string    `json:"desc" firestore:"desc"`
	Slug          string    `json:"slug" firestore:"slug"`
	NameKey       string    `json:"nameKey" firestore:"nameKey"`
	Quantity      int       `json:"quantity" firestore:"quantity"`
	CreatedAt     time.Time `json:"created_at" firestore:"createdAt"`
	UpdatedAt     time.Time `json:"updated_at" firestore:"updatedAt"`
}

func (p *Product) ToMap() map[string]any {
	return map[string]any{
		"id":            p.ID,
		"isActive":      p.IsActive,
		"brand":         p.Brand,
		"name":          p.Name,
		"sku":           p.SKU,
		"size":          p.Size,
		"sizeUnit":      p.SizeUnit,
		"packOf":        p.PackOf,
		"hazardous":     p.Hazardous,
		"category":      p.Category,
		"price":         p.Price,
		"purchasePrice": p.PurchasePrice,
		"desc":          p.Desc,
		"slug":          p.Slug,
		"nameKey":       p.NameKey,
		"quantity":      p.Quantity,
		"createdAt":     p.CreatedAt,
		"updatedAt":     p.UpdatedAt,
	}
}

// ToMinimalMap returns a minimal map of product stored in firestore.
// Product objects are very big and can grow in the future, making sure
// they are stored in a minimal amount of data is important to avoid
// big writes to one document
func (p *Product) ToMinimalMap() map[string]any {
	return map[string]any{
		"id":       p.ID,
		"quantity": p.Quantity,
		"price":    p.Price,
	}
}

/* Setters */
func (p *Product) SetID(id string) {
	p.ID = id
}

func (p *Product) SetIsActive(isActive bool) {
	p.IsActive = isActive
}

func (p *Product) SetBrand(brand string) {
	p.Brand = brand
}

func (p *Product) SetName(name string) {
	p.Name = name
}

func (p *Product) SetSKU(sku string) {
	p.SKU = sku
}

func (p *Product) SetSize(size float64) {
	p.Size = size
}

func (p *Product) SetSizeUnit(sizeUnit string) {
	p.SizeUnit = sizeUnit
}

func (p *Product) SetPackOf(packOf int) {
	p.PackOf = packOf
}

func (p *Product) SetHazardous(hazardous bool) {
	p.Hazardous = hazardous
}

func (p *Product) SetCategory(category string) {
	p.Category = category
}

func (p *Product) SetPrice(price float64) {
	p.Price = price
}

func (p *Product) SetPurchasePrice(purchasePrice float64) {
	p.PurchasePrice = purchasePrice
}

func (p *Product) SetDesc(desc string) {
	p.Desc = desc
}

func (p *Product) SetSlug(slug string) {
	p.Slug = slug
}

func (p *Product) SetNameKey(nameKey string) {
	p.NameKey = nameKey
}

func (p *Product) SetQuantity(quantity int) {
	p.Quantity = quantity
}

func (p *Product) SetCreatedAt(createdAt time.Time) {
	p.CreatedAt = createdAt
}

func (p *Product) SetUpdatedAt(updatedAt time.Time) {
	p.UpdatedAt = updatedAt
}

//Calculate methods

func (p *Product) GetTotalPrice() float64 {
	return p.Price * float64(p.Quantity)
}
func (p *Product) GetTotalPurchasePrice() float64 {
	return p.PurchasePrice * float64(p.Quantity)
}
func (p *Product) GetTotalRevenue() float64 {
	return p.GetTotalPrice() - p.GetTotalPurchasePrice()
}
func (p *Product) GetCorrectWeightInGallons() float64 {
	unit := strings.ToUpper(p.SizeUnit)
	switch unit {
	case "OZ", "OUNCE", "OUNCES":
		return (p.Size / 128) * float64(p.Quantity)
	case "LB", "LBS", "POUND", "POUNDS":
		return (p.Size * 0.125) * float64(p.Quantity)
	case "QT", "QUART", "QUARTS":
		return (p.Size * 0.25) * float64(p.Quantity)
	case "GAL", "GALLON", "GALLONS":
		return p.Size * float64(p.Quantity)
	case "SHEETS": //Fabric softener sheets. Its a rough estimate
		return float64(p.Quantity) * 0.15 / 128
	default:
		return p.Size * float64(p.Quantity)
	}
}

func AreEqualPrices(a, b []Product) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Price != b[i].Price {
			return false
		}
	}
	return true
}

func AreEqualQuantities(a, b []Product) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Quantity != b[i].Quantity {
			return false
		}
	}
	return true
}

//Format methods

func (p *Product) GetFormattedDescription() string {
	return fmt.Sprintf("%s - %s %.2f %s (Pack of %d)", p.Brand, p.Name, p.Size, strings.ToLower(p.SizeUnit), p.PackOf)
}

func (p *Product) GetFormattedUnitPrice() string {
	return fmt.Sprintf("$%.2f", p.Price)
}

func (p *Product) GetFormattedPurchasePrice() string {
	return fmt.Sprintf("$%.2f", p.PurchasePrice)
}

func (p *Product) GetFormattedTotalPrice() string {
	return fmt.Sprintf("$%.2f", p.GetTotalPrice())
}
func (p *Product) GetFormattedTotalPurchasePrice() string {
	return fmt.Sprintf("$%.2f", p.GetTotalPurchasePrice())
}
func (p *Product) GetFormattedTotalRevenue() string {
	return fmt.Sprintf("$%.2f", p.GetTotalRevenue())
}
func (p *Product) GetFormattedQuantity() string {
	return fmt.Sprintf("%d", p.Quantity)
}

func (p *Product) GetFormattedIsHazardous() string {
	if p.Hazardous {
		return "Yes"
	}
	return "No"
}

func (p *Product) GetFormattedTotalWeight() string {
	return fmt.Sprintf("%.2f gal", p.GetCorrectWeightInGallons())
}

func (p *Product) GetFormattedTotalHazardousWeight() string {
	if p.Hazardous {
		return fmt.Sprintf("%.2f gal", p.GetCorrectWeightInGallons())
	}
	return "N/A"
}

func (p *Product) GetFormattedTotalNonHazardousWeight() string {
	if !p.Hazardous {
		return fmt.Sprintf("%.2f gal", p.GetCorrectWeightInGallons())
	}
	return "N/A"
}
