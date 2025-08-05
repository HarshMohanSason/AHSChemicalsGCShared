package qbmodels

import (
	"encoding/json"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
)

type Invoice struct {
	ID          string               `json:"Id,omitempty"`
	SyncToken   string               `json:"SyncToken,omitempty"`
	CustomerRef Reference            `json:"CustomerRef"` // Required
	TxnDate     string               `json:"TxnDate,omitempty"`
	DueDate     string               `json:"DueDate,omitempty"`
	TotalAmt    float64              `json:"TotalAmt,omitempty"`
	Balance     float64              `json:"Balance,omitempty"`
	CurrencyRef *Reference           `json:"CurrencyRef,omitempty"`
	Line        []InvoiceLine        `json:"Line"`
	CustomField []CustomField        `json:"CustomField,omitempty"`
	MetaData    *quickbooks.MetaData `json:"MetaData,omitempty"`
}

func NewInvoice(order *models.Order) *Invoice {

	invoice :=  &Invoice{
		CustomerRef: Reference{Value: order.Customer.ID, Name: order.Customer.Name},
		TxnDate:     time.Now().Format("2006-01-02"),
		DueDate:     time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
		TotalAmt:    order.Total,
	}
	invoice.AddLines(order)
	return invoice
}

func (i *Invoice) ToBytes() ([]byte, error){
	return json.Marshal(i)
}

func (i *Invoice) AddLines(order *models.Order) {
	invoiceLines := make([]InvoiceLine, 0)
	for _ , item := range order.Items{
		line := InvoiceLine{
			DetailType: "SalesItemLineDetail",
			Description: item.GetFormattedDescription(),
			Amount:      item.GetTotalPrice(),
		}
		invoiceLines = append(invoiceLines, line)
	}
	i.Line = invoiceLines
}

type Reference struct {
	Value string `json:"value"`
	Name  string `json:"name,omitempty"`
}

type InvoiceLine struct {
	DetailType          string               `json:"DetailType"`
	Description         string               `json:"Description,omitempty"`
	Amount              float64              `json:"Amount"`
	SalesItemLineDetail *SalesItemLineDetail `json:"SalesItemLineDetail,omitempty"`
	DiscountLineDetail  *DiscountLineDetail  `json:"DiscountLineDetail,omitempty"`
	GroupLineDetail     *GroupLineDetail     `json:"GroupLineDetail,omitempty"`
}

type SalesItemLineDetail struct {
	ItemRef    Reference  `json:"ItemRef"` // Required detail for sales item
	Qty        float64    `json:"Qty,omitempty"`
	UnitPrice  float64    `json:"UnitPrice,omitempty"`
	TaxCodeRef *Reference `json:"TaxCodeRef,omitempty"`
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
	QueryResponse struct {
		Invoice []Invoice `json:"Invoice"`
	} `json:"QueryResponse"`
}

func (qb *QBInvoiceResponse) GetInvoiceID() string{
	return qb.QueryResponse.Invoice[0].ID
}
