package qbmodels

import (
	"encoding/json"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
)

type Invoice struct {
	ID                    string               `json:"Id,omitempty"`
	SyncToken             string               `json:"SyncToken,omitempty"`
	CustomerRef           Reference            `json:"CustomerRef"` // Required
	TxnDate               string               `json:"TxnDate,omitempty"`
	DueDate               string               `json:"DueDate,omitempty"`
	DocNumber             string               `json:"DocNumber,omitempty"`
	PrivateNote           string               `json:"PrivateNote,omitempty"`
	Line                  []Line        	   `json:"Line"` // Required
	CustomField           []CustomField        `json:"CustomField,omitempty"`
	TotalAmt              float64              `json:"TotalAmt,omitempty"`
	Balance               float64              `json:"Balance,omitempty"`
	ApplyTaxAfterDiscount bool                 `json:"ApplyTaxAfterDiscount,omitempty"`
	PrintStatus           string               `json:"PrintStatus,omitempty"`
	EmailStatus           string               `json:"EmailStatus,omitempty"`
	CurrencyRef           *Reference           `json:"CurrencyRef,omitempty"`
	ExchangeRate          float64              `json:"ExchangeRate,omitempty"`
	MetaData              *quickbooks.MetaData `json:"MetaData,omitempty"`
	BillAddr              *QBCustomerAddress   `json:"BillAddr,omitempty"`
	ShipAddr              *QBCustomerAddress   `json:"ShipAddr,omitempty"`
	SalesTermRef          *Reference           `json:"SalesTermRef,omitempty"`
	DueDateBasedOn        string               `json:"DueDateBasedOn,omitempty"`
	TrackingNum           string               `json:"TrackingNum,omitempty"`
	ShipDate              string               `json:"ShipDate,omitempty"`
	ShippingMethodRef     *Reference           `json:"ShippingMethodRef,omitempty"`
	ClassRef              *Reference           `json:"ClassRef,omitempty"`
	DepartmentRef         *Reference           `json:"DepartmentRef,omitempty"`
	LocationRef           *Reference           `json:"LocationRef,omitempty"`
	ARAccountRef          *Reference           `json:"ARAccountRef,omitempty"`
	ProjectRef            *Reference           `json:"ProjectRef,omitempty"`
}

func NewInvoice(order *models.Order) *Invoice {
	invoice := &Invoice{
		CustomerRef:      Reference{Value: order.Customer.ID, Name: order.Customer.Name},
		TxnDate:          order.UpdatedAt.Format("2006-01-02"),
		DueDate:          order.UpdatedAt.AddDate(0, 0, 30).Format("2006-01-02"),
		TotalAmt:         order.Total,
	}
	invoice.AddLines(order)
	return invoice
}

func (i *Invoice) ToBytes() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Invoice) GetDocNumber() string {
	return i.DocNumber
}

func (i *Invoice) AddLines(order *models.Order) {
	invoiceLines := make([]Line, 0)
	for _, item := range order.Items {
		line := Line{
			DetailType:  "SalesItemLineDetail",
			Description: item.GetFormattedDescription(),
			Amount:      item.GetTotalPrice(),
		}
		line.SetSalesItemLineDetail(item)
		invoiceLines = append(invoiceLines, line)
	}
	i.Line = invoiceLines
}

type Reference struct {
	Value string `json:"value"`
	Name  string `json:"name,omitempty"`
}

type Line struct {
	DetailType          string               `json:"DetailType"`
	Description         string               `json:"Description,omitempty"`
	Amount              float64              `json:"Amount"`
	SalesItemLineDetail *SalesItemLineDetail `json:"SalesItemLineDetail,omitempty"`
	DiscountLineDetail  *DiscountLineDetail  `json:"DiscountLineDetail,omitempty"`
	GroupLineDetail     *GroupLineDetail     `json:"GroupLineDetail,omitempty"`
}

func (i *Line) SetSalesItemLineDetail(item *models.Product) {
	detail := &SalesItemLineDetail{
		ItemRef: Reference{
			Value: item.ID,
			Name:  item.Name,
		},
		Qty:       float64(item.Quantity),
		UnitPrice: item.Price,
		TaxCodeRef: &Reference{Value: "TAX"},
	}
	i.SalesItemLineDetail = detail
}

type SalesItemLineDetail struct {
	ItemRef    Reference  `json:"ItemRef"` // Required detail for sales item
	Qty        float64    `json:"Qty,omitempty"`
	UnitPrice  float64    `json:"UnitPrice,omitempty"`
	TaxCodeRef *Reference `json:"TaxCodeRef,omitempty"`
}

type TxnTaxDetail struct {
	TxnTaxCodeRef Reference `json:"TxnTaxCodeRef"`
	TotalTax      float64   `json:"TotalTax"`
	TaxLine       []TaxLine `json:"TaxLine,omitempty"`
}

type TaxLine struct {
	Amount        float64       `json:"Amount"`
	DetailType    string        `json:"DetailType"`
	TaxLineDetail TaxLineDetail `json:"TaxLineDetail"`
}

type TaxLineDetail struct {
	TaxRateRef       Reference `json:"TaxRateRef"`
	NetAmountTaxable float64   `json:"NetAmountTaxable"`
	PercentBased     bool      `json:"PercentBased"`
	TaxPercent       float64   `json:"TaxPercent"`
	TaxInclusiveAmt  float64   `json:"TaxInclusiveAmount,omitempty"`
}

type DiscountLineDetail struct {
	PercentBased    bool    `json:"PercentBased"`
	DiscountPercent float64 `json:"DiscountPercent,omitempty"`
}

type GroupLineDetail struct {
	GroupItemRef Reference `json:"GroupItemRef"`
	Quantity     float64   `json:"Quantity,omitempty"`
}

type CustomField struct {
	DefinitionID string `json:"DefinitionId,omitempty"`
	Name         string `json:"Name,omitempty"`
	Type         string `json:"Type,omitempty"`
	StringValue  string `json:"StringValue,omitempty"`
}

type QBInvoiceResponse struct {
	Invoice Invoice `json:"Invoice"`
	Time    string  `json:"time"`
}

func (qb *QBInvoiceResponse) GetDocNumber() string {
	return qb.Invoice.GetDocNumber()
}