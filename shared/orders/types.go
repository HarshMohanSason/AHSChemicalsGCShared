package orders

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/customers"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/products"
)

type Order struct {
	ID                  string             `json:"ID"`
	Customer            customers.Customer `json:"Customer"`
	Uid                 string             `json:"Uid"`
	SpecialInstructions string             `json:"SpecialInstructions"`
	Items               []products.Item    `json:"Items"`
	TaxRate             float64            `json:"TaxRate"`
	TaxAmount           float64            `json:"TaxAmount"`
	SubTotal            float64            `json:"SubTotal"`
	Total               float64            `json:"Total"`
	Status              string             `json:"Status"`
	PurchaseOrderUrl    string             `json:"PurchaseOrderURL"`
	ShippingManifest    string             `json:"ShippingManifestURL"`
	Invoice             string             `json:"InvoiceURL"`
	CreatedAt           time.Time          `json:"CreatedAt"`
	UpdatedAt           time.Time          `json:"UpdatedAt"`
}
