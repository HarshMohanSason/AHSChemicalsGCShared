package tests

import (
	"context"
	"testing"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/orders"
)

func TestOrderUtils(t *testing.T) {

	order := orders.Order{
		CustomerId:          "1",
		SpecialInstructions: "Handle with care",
		TaxRate:             0.1,
		Items: []orders.OrderItem{
			{
				ProductId: "1",
				Quantity:  2,
				Price:     1.0,
				Size:      1.0,
			},
		},
	}

	err := orders.ValidateOrder(order)
	if err != nil{
		t.Error(err)
	}

	err = orders.CreateOrder(&order)
	if err != nil{
		t.Error(err)
	}

	err = orders.CreateOrderInFirestore(&order, context.Background())
	if err != nil{
		t.Error(err)
	}
	
	t.Logf("Created order is %+v", order)
}
	