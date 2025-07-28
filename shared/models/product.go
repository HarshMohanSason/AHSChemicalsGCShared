package models

import (
	"fmt"
	"strings"
	"time"
)

type ItemMinimal struct {
	ID       string  `firestore:"id"`
	Price    float64 `firestore:"price"`
	Quantity int     `firestore:"quantity"`
}

type Product struct {
	ID        string    `json:"id" firestore:"id"`
	IsActive  bool      `json:"isActive" firestore:"isActive"`
	Brand     string    `json:"brand" firestore:"brand"`
	Name      string    `json:"name" firestore:"name"`
	SKU       string    `json:"sku" firestore:"sku"`
	Size      float64   `json:"size" firestore:"size"`
	SizeUnit  string    `json:"sizeUnit" firestore:"sizeUnit"`
	PackOf    int       `json:"packOf" firestore:"packOf"`
	Hazardous bool      `json:"hazardous" firestore:"hazardous"`
	Category  string    `json:"category" firestore:"category"`
	Price     float64   `json:"price" firestore:"price"`
	Desc      string    `json:"desc" firestore:"desc"`
	Slug      string    `json:"slug" firestore:"slug"`
	NameKey   string    `json:"nameKey" firestore:"nameKey"`
	Quantity  int       `json:"quantity" firestore:"quantity"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
}

func (p *Product) MapToFirestore() map[string]interface{}{
	return map[string]interface{}{
		"id":        p.ID,
		"isActive":  p.IsActive,
		"brand":     p.Brand,
		"name":      p.Name,
		"sku":       p.SKU,
		"size":      p.Size,
		"sizeUnit":  p.SizeUnit,
		"packOf":    p.PackOf,
		"hazardous": p.Hazardous,
		"category":  p.Category,
		"price":     p.Price,
		"desc":      p.Desc,
		"slug":      p.Slug,
		"nameKey":   p.NameKey,
		"quantity":  p.Quantity,
		"createdAt": p.CreatedAt,
		"updatedAt": p.UpdatedAt,
	}
}

//Format methods

func (p *Product) FormatProductDisplay() string {
	return fmt.Sprintf("%s - %s %.2f %s (Pack of %d)", p.Brand, p.Name, p.Size, p.SizeUnit, p.PackOf)
}

func (p *Product) FormatUnitPrice() string {
	return fmt.Sprintf("$%.2f", p.Price)
}

func (p *Product) FormatTotalPrice() string {
	return fmt.Sprintf("$%.2f", p.GetTotalPrice())
}

func (p *Product) FormatQuantity() string {
	return fmt.Sprintf("%d", p.Quantity)
}

func (p *Product) FormatIsHazardous() string {
	if p.Hazardous {
		return "Yes"
	}
	return "No"
}

func (p *Product) FormatNetWeight() string {
	return fmt.Sprintf("%.2f gal", p.GetCorrectWeightInGallons())
}

func (p *Product) FormatHazardousWeight() string {
	if p.Hazardous {
		return fmt.Sprintf("%.2f gal", p.GetCorrectWeightInGallons())
	}
	return "N/A"
}

func (p *Product) FormatNonHazardousWeight() string {
	if !p.Hazardous {
		return fmt.Sprintf("%.2f gal",  p.GetCorrectWeightInGallons())
	}
	return "N/A"
}

//Calculate methods

func (p *Product) GetTotalPrice() float64 {
	return p.Price * float64(p.Quantity)
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