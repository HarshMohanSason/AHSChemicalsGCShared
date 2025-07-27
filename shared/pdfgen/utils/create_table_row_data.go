package utils

import "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"

// PURCHASE ORDER PDF TABLE ROWS GENERATION
func CreateTableRowValuesForPurchaseOrderPDF(order *models.Order) [][]string {
	var mappedOrders = make([][]string, 0)

	for _, item := range order.Items {
		items := []string{item.SKU, item.FormatProductDisplay(), item.FormatQuantity(), item.FormatUnitPrice(), item.FormatTotalPrice()}
		mappedOrders = append(mappedOrders, items)
	}
	return mappedOrders
}

// SHIPPING MANIFEST PDF TABLE ROWS GENERATION
func CreateTableRowValuesForShippingManifestPDF(order *models.Order) [][]string {
	var mappedOrders = make([][]string, 0)
	const (
		typeContainer = "Carton"
		productClass  = "55.0"
	)
	for _, item := range order.Items {
		items := []string{item.FormatQuantity(), item.FormatIsHazardous(), typeContainer, item.FormatProductDisplay(), productClass, item.SKU, item.FormatNetWeight(), item.FormatNonHazardousWeight(), item.FormatHazardousWeight()}
		mappedOrders = append(mappedOrders, items)
	}
	return mappedOrders
}

// INVOICE PDF TABLE ROWS GENERATION 
func CreateTableRowValuesForInvoicePDF(order *models.Order) [][]string {
	var mappedOrders = make([][]string, 0)
	for _, item := range order.Items {
		items := []string{item.FormatProductDisplay(), item.FormatQuantity(), item.FormatUnitPrice(), item.FormatTotalPrice()}
		mappedOrders = append(mappedOrders, items)
	}
	return mappedOrders
}