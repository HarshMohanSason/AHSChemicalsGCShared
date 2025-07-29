package tests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/repositories"
)

func TestSyncProductsToFirestore(t *testing.T) {

	jsonData := `{
		"QueryResponse": {
			"Item": [
				{
					"Id": "1",
					"Name": "Coffee Mug",
					"Type": "Inventory",
					"UnitPrice": 9.99
				},
				{
					"Id": "2",
					"Name": "T-shirt",
					"Type": "NonInventory",
					"UnitPrice": 14.99
				}
			]
		}
	}`
	var qbItemResponse qbmodels.QBItemsResponse
	err := json.Unmarshal([]byte(jsonData), &qbItemResponse)
	if err != nil {
		t.Error(err)
	}

	err = repositories.SyncQuickbookProductRespToFirestore(&qbItemResponse, context.Background())

	if err != nil {
		t.Error(err)
	}
}

func TestFetchAllProductsFromFirestore(t *testing.T) {
	products, err := repositories.FetchAllProductsFromFirestore(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log("Length of fetched products is: ", len(products))
}