package quickbooks

import "github.com/HarshMohanSason/AHSChemicalsGCShared/shared"

type Item struct {
	Id                string            `json:"Id"`
	Name              string            `json:"Name"`
	Description       string            `json:"Description,omitempty"`
	UnitPrice         float64           `json:"UnitPrice,omitempty"`
	Type              string            `json:"Type,omitempty"`
	Active            bool              `json:"Active"`
	TrackQtyOnHand    bool              `json:"TrackQtyOnHand,omitempty"`
	QtyOnHand         float64           `json:"QtyOnHand,omitempty"`
	InvStartDate      string            `json:"InvStartDate,omitempty"`
	IncomeAccountRef  *shared.Reference `json:"IncomeAccountRef,omitempty"`
	ExpenseAccountRef *shared.Reference `json:"ExpenseAccountRef,omitempty"`
	AssetAccountRef   *shared.Reference `json:"AssetAccountRef,omitempty"`
	PurchaseDesc      string            `json:"PurchaseDesc,omitempty"`
	PurchaseCost      float64           `json:"PurchaseCost,omitempty"`
	Taxable           bool              `json:"Taxable,omitempty"`
	SKU               string            `json:"Sku,omitempty"`
	SyncToken         string            `json:"SyncToken,omitempty"`
	ParentRef         *shared.Reference `json:"ParentRef,omitempty"` // Category reference
	MetaData          *shared.MetaData  `json:"MetaData,omitempty"`
}
