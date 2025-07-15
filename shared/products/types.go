package products

import "github.com/HarshMohanSason/AHSChemicalsGCShared/shared"

type Ref struct {
	Value string `json:"value"`
	Name  string `json:"name,omitempty"`
}

type SalesOrPurchase struct {
	Desc       string  `json:"Desc,omitempty"`
	Price      float64 `json:"Price,omitempty"`
	AccountRef *Ref    `json:"AccountRef,omitempty"`
	TaxCodeRef *Ref    `json:"TaxCodeRef,omitempty"`
}

// Note** The Quantity field is different from the QtyOnHand. Quantity represents how much of the item was ordered when the order was placed. QtyOnHand is the quantity of the item on hand which is received from quickbooks api when sycning.
// Size, SizeUnit, PackOf, Slug and Brand are not present in quickbooks api response. These are custom fields created when syncing. 
type Item struct {
	ID                  string           `json:"Id,omitempty"`
	Name                string           `json:"Name"`
	Quantity            int              `json:"Quantity"`
	Brand               string           `json:"Brand"`
	SKU                 string           `json:"SKU"`
	Size                float64          `json:"Size"`
	SizeUnit            string           `json:"SizeUnit"`
	PackOf              int              `json:"PackOf"`
	Description         string           `json:"Description,omitempty"`
	Slug                string           `json:"Slug"`
	NameKey             string           `json:"NameKey"`
	Active              bool             `json:"Active"`
	FullyQualifiedName  string           `json:"FullyQualifiedName,omitempty"`
	Taxable             bool             `json:"Taxable"`
	UnitPrice           float64          `json:"UnitPrice,omitempty"`
	Type                string           `json:"Type"` // "Inventory", "NonInventory", "Service"
	IncomeAccountRef    *Ref             `json:"IncomeAccountRef,omitempty"`
	PurchaseDesc        string           `json:"PurchaseDesc,omitempty"`
	PurchaseCost        float64          `json:"PurchaseCost,omitempty"`
	ExpenseAccountRef   *Ref             `json:"ExpenseAccountRef,omitempty"`
	TrackQtyOnHand      bool             `json:"TrackQtyOnHand,omitempty"`
	QtyOnHand           float64          `json:"QtyOnHand,omitempty"`
	InvStartDate        string           `json:"InvStartDate,omitempty"`
	AssetAccountRef     *Ref             `json:"AssetAccountRef,omitempty"`
	ParentRef           *Ref             `json:"ParentRef,omitempty"`
	Level               int              `json:"Level,omitempty"`
	SalesTaxCodeRef     *Ref             `json:"SalesTaxCodeRef,omitempty"`
	PurchaseTaxCodeRef  *Ref             `json:"PurchaseTaxCodeRef,omitempty"`
	UnitOfMeasure       string           `json:"UnitOfMeasure,omitempty"`
	PreferredVendorRef  *Ref             `json:"PreferredVendorRef,omitempty"`
	TrackQtyAndValue    bool             `json:"TrackQtyAndValue,omitempty"`
	MetaData            *shared.MetaData `json:"MetaData,omitempty"`
	SalesOrPurchaseInfo *SalesOrPurchase `json:"SalesOrPurchase,omitempty"`
}
