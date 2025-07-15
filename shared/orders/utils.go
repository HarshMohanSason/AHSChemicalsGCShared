package orders

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/products"
)

func ValidateOrderItems(item *products.Item) error {
	if item.ID == "" {
		return errors.New("validation failed: product ID is required")
	}
	if item.Quantity == 0 {
		return fmt.Errorf("validation failed: quantity is required for product ID '%s'", item.ID)
	}
	if item.UnitPrice == 0 {
		return fmt.Errorf("validation failed: price is required for product ID '%s'", item.ID)
	}
	if item.Size == 0 {
		return fmt.Errorf("validation failed: size is required for product ID '%s'", item.ID)
	}
	return nil
}

func ValidateOrder(order *Order) error {
	if order.Customer.ID == "" {
		return errors.New("validation failed: customer ID is required")
	}
	if len(order.Items) == 0 {
		return errors.New("validation failed: order must contain at least one item")
	}
	if order.TaxRate == 0 {
		return errors.New("validation failed: tax rate must be provided")
	}
	for _, item := range order.Items {
		if err := ValidateOrderItems(&item); err != nil {
			return err
		}
	}
	return nil
}

func GetCorrectOrderPricesForOrder(order *Order) error {
	productPricesMap := make(map[string]float64)
	for _, item := range order.Items {
		productPricesMap[item.ID] = item.UnitPrice
	}

	err := FetchCustomerPriceForEachProductID(productPricesMap, order.Customer.ID, context.Background())
	if err != nil {
			return fmt.Errorf("failed to fetch customer-specific pricing: %w", err)
	}

	for i, item := range order.Items {
		order.Items[i].UnitPrice = productPricesMap[item.ID]
	}
	return nil
}

func CalculateSubTotal(order *Order) float64 {
	subTotal := 0.0
	for _, item := range order.Items {
		subTotal += float64(item.Quantity) * item.UnitPrice
	}
	return subTotal
}

func CreateOrder(order *Order) error {
	err := GetCorrectOrderPricesForOrder(order)
	if err != nil {
		return err
	}

	order.SubTotal = CalculateSubTotal(order)
	order.TaxAmount = order.TaxRate * order.SubTotal
	order.Total = order.SubTotal + order.TaxAmount
	order.Status = "PENDING"

	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	return nil
}

func FormatOrderForFirestore(order *Order) map[string]any {
	formattedItems := make([]map[string]any, 0)
	for _, item := range order.Items {
		formattedItems = append(formattedItems, map[string]any{
			"Id":        item.ID,
			"Quantity":  item.Quantity,
		})
	}

	return map[string]any{
		"CustomerId":          order.Customer.ID,
		"Uid":                 order.Uid,
		"SpecialInstructions": order.SpecialInstructions,
		"Items":               formattedItems,
		"TaxRate":             order.TaxRate,
		"TaxAmount":           order.TaxAmount,
		"SubTotal":            order.SubTotal,
		"Total":               order.Total,
		"Status":              order.Status,
		"PurchaseOrderURL":    order.PurchaseOrderUrl,
		"CreatedAt":           order.CreatedAt,
		"UpdatedAt":           order.UpdatedAt,
	}
}

func CreateOrderForPurchaseOrderPDF(order *Order) [][]string {
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