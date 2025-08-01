package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/constants"
)

// Order struct represents an order placed by a user.
type Order struct {
	ID                  string    `json:"id" firestore:"omitempty"`
	Customer            Customer  `json:"customer" firestore:"customer"`
	Uid                 string    `json:"uid" firestore:"uid"` // User ID of placed the order
	SpecialInstructions string    `json:"specialInstructions" firestore:"specialInstructions"`
	Items               []Product `json:"items" firestore:"items"`
	TaxRate             float64   `json:"taxRate" firestore:"taxRate"`
	TaxAmount           float64   `json:"taxAmount" firestore:"taxAmount"`
	SubTotal            float64   `json:"subTotal" firestore:"subTotal"`
	Total               float64   `json:"total" firestore:"total"`
	Status              string    `json:"status" firestore:"status"`
	CreatedAt           time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt" firestore:"updatedAt"`
}

// CreateCompleteOrder creates the complete order
// Calculates the subtotal, tax amount and total and adds a default status pending.
//
// Param:
//   - correctPrices map[string]float64 prices mapped with their corresponding item ID
func (o *Order) CreateCompleteOrder(correctPrices map[string]float64) {
	if correctPrices == nil {
		return
	}
	o.SetItemPrices(correctPrices)
	o.calcSubtotal()
	o.calcTaxAmount()
	o.calcTotal()
	o.SetStatus(constants.OrderStatusPending)
	o.setCreatedAt()
	o.SetUpdatedAt()
}

/* Setters */

func (o *Order) SetID(id string) {
	o.ID = id
}
func (o *Order) SetUID(uid string){
	o.Uid = uid
}
func (o *Order) SetStatus(status string) {
	o.Status = status
}

func (o *Order) SetItemPrices(correctPrices map[string]float64) {
	for _, item := range o.Items {
		item.SetPrice(correctPrices[item.ID])
	}
}

func (o *Order) setCreatedAt() {
	o.CreatedAt = time.Now()
}

func (o *Order) SetUpdatedAt() {
	o.CreatedAt = time.Now()
}

/* Calculations */

func (o *Order) calcSubtotal() {
	for _, item := range o.Items {
		o.SubTotal += item.GetTotalPrice()
	}
}

func (o *Order) calcTaxAmount() {
	o.TaxAmount = o.TaxRate * o.SubTotal
}

func (o *Order) calcTotal() {
	o.Total = o.SubTotal + o.TaxAmount
}

/* Converters */

func (o *Order) ToMap() map[string]any {
	return map[string]any{
		"id":                  o.ID,
		"customerId":          o.Customer.ID,
		"customerName":        strings.ToLower(o.Customer.Name),
		"uid":                 o.Uid,
		"specialInstructions": o.SpecialInstructions,
		"items":               o.toMapItems(),
		"taxRate":             o.TaxRate,
		"taxAmount":           o.TaxAmount,
		"subTotal":            o.SubTotal,
		"total":               o.Total,
		"status":              o.Status,
		"createdAt":           o.CreatedAt,
		"updatedAt":           o.UpdatedAt,
	}
}

func (o *Order) toMapItems() []map[string]any {
	minimalItems := make([]map[string]any, 0)
	for _, items := range o.Items {
		minimalItems = append(minimalItems, items.ToMinimalMap())
	}
	return minimalItems
}

// Returns an array of product IDs. Comes in handy for bulk firestore operations
func (o *Order) ToProductIDs() []string {
	productIDs := make([]string, 0)
	for _, item := range o.Items {
		productIDs = append(productIDs, item.ID)
	}
	return productIDs
}

// This is used to convert the order items array stored in firestore to a
// complete order object.
// Note: the products map[string]Product should be fetched from firestore which
// contains the original products mapped with their id's.
func (o *Order) ToCompleteOrderItemsFromMinimal(products map[string]Product) {
	if products == nil {
		return
	}
	for i, item := range o.Items {
		o.Items[i].SetIsActive(products[item.ID].IsActive)
		o.Items[i].SetBrand(products[item.ID].Brand)
		o.Items[i].SetName(products[item.ID].Name)
		o.Items[i].SetSKU(products[item.ID].SKU)
		o.Items[i].SetSize(products[item.ID].Size)
		o.Items[i].SetSizeUnit(products[item.ID].SizeUnit)
		o.Items[i].SetPackOf(products[item.ID].PackOf)
		o.Items[i].SetHazardous(products[item.ID].Hazardous)
		o.Items[i].SetCategory(products[item.ID].Category)
		o.Items[i].SetDesc(products[item.ID].Desc)
		o.Items[i].SetSlug(products[item.ID].Slug)
		o.Items[i].SetNameKey(products[item.ID].NameKey)
		o.Items[i].SetCreatedAt(products[item.ID].CreatedAt)
		o.Items[i].SetUpdatedAt(products[item.ID].UpdatedAt)
	}
}

/* Formatters */

func (o *Order) GetFormattedTotal() string {
	return fmt.Sprintf("$%.2f", o.Total)
}

func (o *Order) GetFormattedSubTotal() string {
	return fmt.Sprintf("$%.2f", o.SubTotal)
}

func (o *Order) GetFormattedTaxAmount() string {
	return fmt.Sprintf("$%.2f", o.TaxAmount)
}

func (o *Order) GetFormattedTaxRate() string {
	return fmt.Sprintf("%.2f%%", o.TaxRate)
}

func (o *Order) GetFormattedTotalItems() string {
	totalUnits := 0
	for _, item := range o.Items {
		totalUnits += item.Quantity
	}
	return strconv.Itoa(totalUnits)
}

func (o *Order) GetFormattedNetWeight() string {
	weight := 0.0
	for _, item := range o.Items {
		weight += item.GetCorrectWeightInGallons()
	}
	return fmt.Sprintf("%.2f gal", weight)
}

func (o *Order) GetFormattedNetNonHazardousWeight() string {
	weight := 0.0
	for _, item := range o.Items {
		if !item.Hazardous {
			weight += item.GetCorrectWeightInGallons()
		}
	}
	return fmt.Sprintf("%.2f gal", weight)
}

func (o *Order) GetFormattedNetHazardousWeight() string {
	weight := 0.0
	for _, item := range o.Items {
		if item.Hazardous {
			weight += item.GetCorrectWeightInGallons()
		}
	}
	return fmt.Sprintf("%.2f gal", weight)
}
