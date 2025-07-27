package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Order struct represents an order placed by a user.
type Order struct {
	ID                  string        `json:"id" firestore:"omitempty"`
	Customer            Customer      `json:"customer" firestore:"customer"`
	Uid                 string        `json:"uid" firestore:"uid"` // User ID who placed the order
	SpecialInstructions string        `json:"specialInstructions" firestore:"specialInstructions"`
	Items               []Product     `json:"fullitems" firestore:"omitempty"`
	MinimalItems        []ItemMinimal `json:"items" firestore:"items"`
	TaxRate             float64       `json:"taxRate" firestore:"taxRate"`
	TaxAmount           float64       `json:"taxAmount" firestore:"taxAmount"`
	SubTotal            float64       `json:"subTotal" firestore:"subTotal"`
	Total               float64       `json:"total" firestore:"total"`
	Status              string        `json:"status" firestore:"status"`
	CreatedAt           time.Time     `json:"createdAt" firestore:"createdAt"`
	UpdatedAt           time.Time     `json:"updatedAt" firestore:"updatedAt"`
}

// Format methods

func (o *Order) FormatTotal() string {
	return fmt.Sprintf("%.2f", o.Total)
}

func (o *Order) FormatSubTotal() string {
	return fmt.Sprintf("%.2f", o.SubTotal)
}

func (o *Order) FormatTaxAmount() string {
	return fmt.Sprintf("%.2f", o.TaxAmount)
}

func (o *Order) FormatTaxRate() string {
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
	return fmt.Sprintf("%.2f", weight)
}

func (o *Order) GetFormattedNonHazardousWeight() string {
	weight := 0.0
	for _, item := range o.Items {
		if !item.Hazardous {
			weight += item.GetCorrectWeightInGallons()
		}
	}
	return fmt.Sprintf("%.2f GAL", weight)
}

func (o *Order) GetFormattedHazardousWeight() string {
	weight := 0.0
	for _, item := range o.Items {
		if item.Hazardous {
			weight += item.GetCorrectWeightInGallons()
		}
	}
	return fmt.Sprintf("%.2f GAL", weight)
}

func (o *Order) CreatePricesMap() map[string]float64 {
	mappedPrices := make(map[string]float64)
	for _, item := range o.Items {
		mappedPrices[item.ID] = item.Price
	}
	return mappedPrices
}

func (o *Order) ToFirestoreMap() map[string]any {
	return map[string]any{
		"id":                  o.ID,
		"customerId":          o.Customer.ID,
		"customerName":        o.Customer.Name,
		"uid":                 o.Uid,
		"specialInstructions": o.SpecialInstructions,
		"items":               o.MinimalItems,
		"taxRate":             o.TaxRate,
		"taxAmount":           o.TaxAmount,
		"subTotal":            o.SubTotal,
		"total":               o.Total,
		"status":              o.Status,
		"createdAt":           o.CreatedAt,
		"updatedAt":           o.UpdatedAt,
	}
}

// Calculations

func (o *Order) CalcTaxAmount() {
	o.TaxAmount = o.TaxRate * o.SubTotal
}

func (o *Order) CalcSubtotal() {
	for _, item := range o.Items {
		o.SubTotal += item.GetTotalPrice()
	}
}

func (o *Order) CalcTotal() {
	o.Total = o.SubTotal + o.TaxAmount
}

// Validate does basic and empty check to make sure the necessary fields are set
func (o *Order) Validate() error {
	if o.Uid == "" {
		return errors.New("No user id found when order was placed")
	}
	if len(o.MinimalItems) == 0 {
		return errors.New("No items found in order")
	}
	if o.TaxRate == 0 {
		return errors.New("No tax rate found in order")
	}
	if o.Customer.ID == "" {
		return errors.New("No customer found for this order")
	}
	if o.Customer.Name == "" {
		return errors.New("No customer name found for this order")
	}
	return nil
}