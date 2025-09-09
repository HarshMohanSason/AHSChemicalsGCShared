package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/constants"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/utils"
)

type Order struct {
	ID                  string     `json:"id"`
	Customer            *Customer  `json:"customer" firestore:"customer"`
	Uid                 string     `json:"uid" firestore:"uid"` // User ID of placed the order
	SpecialInstructions string     `json:"specialInstructions" firestore:"specialInstructions"`
	Items               []*Product `json:"items" firestore:"items"`
	TaxRate             float64    `json:"taxRate" firestore:"taxRate"`
	TaxAmount           float64    `json:"taxAmount" firestore:"taxAmount"`
	SubTotal            float64    `json:"subTotal" firestore:"subTotal"`
	Total               float64    `json:"total" firestore:"total"`
	Status              string     `json:"status" firestore:"status"`
	CreatedAt           time.Time  `json:"createdAt" firestore:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt" firestore:"updatedAt"`
	TimeZone            string     `json:"timeZone" firestore:"timeZone"`
}

// Small struct to fetch order if only order id is passed form frontend
type OrderIDPaylod struct {
	OrderID string `json:"orderId"`
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
	o.getSubTotal()
	o.getTaxAmount()
	o.getTotal()
	o.SetStatus(constants.OrderStatusPending)
	time := time.Now().UTC()
	o.setCreatedAt(time)
	o.SetUpdatedAt(time)
}

func (o *Order) UpdateOrderBill() {
	o.getSubTotal()
	o.getTaxAmount()
	o.getTotal()
}

//Setters

func (o *Order) SetID(id string) {
	o.ID = id
}

func (o *Order) SetCustomer(customer *Customer) {
	o.Customer = customer
}

func (o *Order) SetUID(uid string) {
	o.Uid = uid
}

func (o *Order) SetStatus(status string) {
	o.Status = status
}

func (o *Order) SetItemPrices(correctPrices map[string]float64) {
	for i, item := range o.Items {
		o.Items[i].SetPrice(correctPrices[item.ID])
	}
}

func (o *Order) setCreatedAt(createdAt time.Time) {
	o.CreatedAt = createdAt
}

func (o *Order) SetUpdatedAt(updatedAt time.Time) {
	o.UpdatedAt = updatedAt
}

//Getters

func (o *Order) getSubTotal() {
	subTotal := 0.0
	for _, item := range o.Items {
		subTotal += item.GetTotalPrice()
	}
	o.SubTotal = subTotal
}

func (o *Order) getTaxAmount() {
	o.TaxAmount = o.TaxRate * o.SubTotal
}

func (o *Order) getTotal() {
	o.Total = o.SubTotal + o.TaxAmount
}

// Gets the total cost of goods with their purchase prices * quantity
func (o *Order) GetTotalCOG() float64 {
	totalCOG := 0.0
	for _, item := range o.Items {
		totalCOG += item.GetTotalPurchasePrice()
	}
	return totalCOG
}

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
	return fmt.Sprintf("%.2f%%", o.TaxRate*100)
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

func (o *Order) GetFormattedCOG() string {
	return fmt.Sprintf("$%.2f", o.GetTotalCOG())
}

// Subtracting from subtotal not total because total includes the tax.
func (o *Order) GetFormattedTotalRevenue() string {
	return fmt.Sprintf("$%.2f", o.SubTotal-o.GetTotalCOG())
}

func (o *Order) GetLocalCreatedAtTime() time.Time {
	localTime, err := utils.ConvertUTCToLocalTimeZoneWithFormat(o.CreatedAt, o.TimeZone)
	if err != nil {
		return o.CreatedAt
	}
	return localTime
}

func (o *Order) GetLocalUpdatedAtTime() time.Time {
	localTime, err := utils.ConvertUTCToLocalTimeZoneWithFormat(o.UpdatedAt, o.TimeZone)
	if err != nil {
		return o.UpdatedAt
	}
	return localTime
}

// Converts the order object to a map that can be stored in firestore.
func (o *Order) ToMap() map[string]any {
	return map[string]any{
		"customerId":          o.Customer.ID,
		"customerName":        strings.ToLower(o.Customer.Name),
		"uid":                 o.Uid,
		"specialInstructions": o.SpecialInstructions,
		"items":               o.ToMapItems(),
		"taxRate":             o.TaxRate,
		"taxAmount":           o.TaxAmount,
		"subTotal":            o.SubTotal,
		"total":               o.Total,
		"status":              o.Status,
		"createdAt":           o.CreatedAt,
		"updatedAt":           o.UpdatedAt,
		"timeZone":            o.TimeZone,
	}
}

// Converts the order items array stored in firestore to a minimal item object.
func (o *Order) ToMapItems() []map[string]any {
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
//
// Note: the products map[string]Product should be fetched from firestore which
// contains the original products mapped with their id's. In addition to this,
// the only thing which is not set in the order is the price of the product. The
// price for each product is always used from the items minimal map which already exists
// in the order. Only time the price is fetched when the order is created.
func (o *Order) ToCompleteOrderItemsFromMinimal(products map[string]*Product) {
	if products == nil {
		return
	}
	for i, item := range o.Items {
		o.Items[i].SetIsActive(products[item.ID].IsActive)
		o.Items[i].SetBrand(products[item.ID].Brand)
		o.Items[i].SetName(products[item.ID].Name)
		o.Items[i].SetSKU(products[item.ID].SKU)
		o.Items[i].SetSize(products[item.ID].Size)
		o.Items[i].SetPurchasePrice(products[item.ID].PurchasePrice)
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

// Converts the order items array to a map of ID:Product
func (o *Order) ToItemMap() map[string]*Product {
	idMap := make(map[string]*Product)
	for _, item := range o.Items {
		idMap[item.ID] = item
	}
	return idMap
}

// Tracks changes when updating a particular order
type TrackOrderChange struct {
	StatusChanged bool
	ItemsChanged  bool
}

func NewOrderTracker() *TrackOrderChange {
	return &TrackOrderChange{
		StatusChanged: false,
		ItemsChanged:  false,
	}
}

func (t *TrackOrderChange) SetStatusChanged(statusChanged bool) {
	t.StatusChanged = statusChanged
}

func (t *TrackOrderChange) SetItemsChanged(new, old []*Product) {
	if !AreEqualPrices(new, old) || !AreEqualQuantities(new, old) {
		t.ItemsChanged = true
	}
}

func (t *TrackOrderChange) HasChanges() bool {
	return t.StatusChanged || t.ItemsChanged
}

func (t *TrackOrderChange) IsOnlyStatusChanged() bool {
	return t.StatusChanged && !t.ItemsChanged
}

func (t *TrackOrderChange) TrackOrderChanges(editedOrder, originalOrder *Order) {
	if editedOrder.Status != originalOrder.Status {
		t.SetStatusChanged(true)
	}
	t.SetItemsChanged(editedOrder.Items, originalOrder.Items)
}
