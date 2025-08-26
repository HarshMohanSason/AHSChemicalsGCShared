// package services contains the business logic for all models from the models package
package services

import (

	"cloud.google.com/go/firestore"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

// GetUpdatedCustomerDetails returns the updated customer details. Used for creating a map
// of what has changed when the customer is updated in quickbooks. This is used by the webhook
// entitiy processor cloud event handler when it receives an update event notificaiton from the
// webhook cloud function. The updated customer is fetched via the cloud function to get a single product
// from quickbooks whereas the old customer is fetched from firestore. Returns a nil map
// if there is no change in the objects.
func GetUpdatedCustomerDetails(updated, oldCustomer *models.Customer) map[string]any {
	if updated == nil || oldCustomer == nil {
		return nil
	}

	changedValues := make(map[string]any)

	if updated.IsActive != oldCustomer.IsActive {
		changedValues["isActive"] = updated.IsActive
	}
	if updated.Name != oldCustomer.Name {
		changedValues["name"] = updated.Name
	}
	if updated.Email != oldCustomer.Email {
		changedValues["email"] = updated.Email
	}
	if updated.Phone != oldCustomer.Phone {
		changedValues["phone"] = updated.Phone
	}
	if updated.Address1 != oldCustomer.Address1 {
		changedValues["address1"] = updated.Address1
	}
	if updated.City != oldCustomer.City {
		changedValues["city"] = updated.City
	}
	if updated.State != oldCustomer.State {
		changedValues["state"] = updated.State
	}
	if updated.Zip != oldCustomer.Zip {
		changedValues["zip"] = updated.Zip
	}
	if updated.Country != oldCustomer.Country {
		changedValues["country"] = updated.Country
	}

	if len(changedValues) == 0 {
		return nil
	}
	changedValues["updatedAt"] = firestore.ServerTimestamp

	return changedValues
}
