package qbmodels

import (
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
)

type QBEstimate struct {
	ID             string               `json:"Id,omitempty"`
	SyncToken      string               `json:"SyncToken,omitempty"`
	MetaData       *quickbooks.MetaData `json:"MetaData,omitempty"`
	CustomerRef    *Reference           `json:"CustomerRef"`
	TotalAmt       float64              `json:"TotalAmt,omitempty"`
	Balance        float64              `json:"Balance,omitempty"`
	TxnDate        string               `json:"TxnDate,omitempty"`
	ExpirationDate string               `json:"ExpirationDate,omitempty"`
	Line           []Line               `json:"Line,omitempty"`
	TxnTaxDetail   *TxnTaxDetail        `json:"TxnTaxDetail,omitempty"`
	CurrencyRef    *Reference           `json:"CurrencyRef,omitempty"`
	CustomField    []CustomField        `json:"CustomField,omitempty"`
}

func NewQBEstimate(order *models.Order) *QBEstimate {
	QBEstimate := &QBEstimate{
		CustomerRef: &Reference{
			Value: order.Customer.ID,
			Name:  order.Customer.Name,
		},
		TotalAmt: order.Total,
		TxnDate:  order.CreatedAt.Format("2006-01-02"),
	}
	QBEstimate.AddLines(order)
	return QBEstimate
}

func (i *QBEstimate) AddLines(order *models.Order) {
	lines := make([]Line, 0)
	for _, item := range order.Items {
		line := Line{
			DetailType:  "SalesItemLineDetail",
			Description: item.GetFormattedDescription(),
			Amount:      item.GetTotalPrice(),
		}
		line.SetSalesItemLineDetail(item)
		lines = append(lines, line)
	}
	i.Line = lines
}

// Wrapper for the api response
type QBEstimateResponse struct {
	Estimate *QBEstimate `json:"Estimate"`
}
