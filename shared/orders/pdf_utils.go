package orders

import (
	"fmt"
	"strconv"
)

// CreateTableRowValuesForPurchaseOrderPDF generates a 2D array of order items formatted for the Purchase Order PDF output.
// Each inner slice represents a row in the PDF table with the following columns:
// [SKU, Description, Quantity, Unit Price, Total Product Price].
//
// The description field is constructed using the item's brand, name, size, unit, and pack size.
// All monetary values are formatted with a dollar sign and two decimal places.
//
// Parameters:
//   - order: Pointer to an Order object containing the items to be converted.
//
// Returns:
//   - [][]string: A 2D string array where each row corresponds to an order item.
func CreateTableRowValuesForPurchaseOrderPDF(order *Order) [][]string {
	var mappedOrders = make([][]string, 0)

	for _, item := range order.Items {
		sku := item.SKU
		description := fmt.Sprintf("%s - %s %.2f %s (Pack of %d)", item.Brand, item.Name, item.Size, item.SizeUnit, item.PackOf)
		quantity := strconv.Itoa(item.Quantity)
		price := fmt.Sprintf("$%.2f", item.UnitPrice)
		total := fmt.Sprintf("$%.2f", item.UnitPrice*float64(item.Quantity))

		items := []string{sku, description, quantity, price, total}
		mappedOrders = append(mappedOrders, items)
	}
	return mappedOrders
}