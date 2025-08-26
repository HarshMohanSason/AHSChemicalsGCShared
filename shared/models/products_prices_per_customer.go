package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type ProductPricePerCustomer struct {
	ProductName string    `json:"productName" firsetore:"productName"`
	Brand       string    `json:"brand" firsetore:"brand"`
	ProductID   string    `json:"productId" firsetore:"productId"`
	CustomerID  string    `json:"customerId" firsetore:"customerId"`
	Price       float64   `json:"price" firsetore:"price"`
	CreatedAt   time.Time `json:"createdAt" firsetore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firsetore:"updatedAt"`
}

func (ppc *ProductPricePerCustomer) ToMap() map[string]any {
	return map[string]any{
		"productName": ppc.ProductName,
		"brand":       ppc.Brand,
		"productId":   ppc.ProductID,
		"customerId":  ppc.CustomerID,
		"price":       ppc.Price,
		"createdAt":   firestore.ServerTimestamp,
		"updatedAt":   firestore.ServerTimestamp,
	}
}

func CreateProductPricePerCustomer(product *Product, customerID string) *ProductPricePerCustomer {
	return &ProductPricePerCustomer{
		ProductName: product.Name,
		Brand:       product.Brand,
		ProductID:   product.ID,
		CustomerID:  customerID,
		Price:       product.Price,
	}
}
