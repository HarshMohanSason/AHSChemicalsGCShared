package services

import (
	"errors"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/constants"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

// approveOrder changes the status of an order to APPROVED
// An order can only be approved if
//   - it is in PENDING state
//   - it is in REJECTED state and no more than 30 days have passed since rejection
//
// Parameters:
//   - editedOrder: Pointer to an Order object containing the edited order.
//   - originalOrder: Pointer to an Order object containing the original order.
//
// Returns:
//   - error: An error object if edited status is incorrect .
func approveOrder(editedOrder, originalOrder *models.Order) error {
	switch originalOrder.Status {
	case constants.OrderStatusPending:
		editedOrder.Status = constants.OrderStatusApproved
		return nil
	case constants.OrderStatusRejected:
		if time.Since(originalOrder.UpdatedAt) > 30*24*time.Hour {
			return errors.New("Cannot approve: more than 30 days have passed since rejection. Place a new order instead")
		}
		editedOrder.Status = constants.OrderStatusApproved
		return nil
	default:
		return errors.New("Order can only be approved if it is in PENDING or REJECTED state")
	}
}

// rejectOrder changes the status of an order to REJECTED.
// An order can only be rejected if
//   - it is in PENDING state
//
// Parameters:
//   - editedOrder: Pointer to an Order object containing the edited order.
//   - originalOrder: Pointer to an Order object containing the original order.
//
// Returns:
//   - error: An error object if edited status is incorrect .
func rejectOrder(editedOrder *models.Order, originalOrderStatus string) error {
	switch originalOrderStatus {
	case constants.OrderStatusPending:
		editedOrder.Status = constants.OrderStatusRejected
		return nil
	default:
		return errors.New("Order can only be rejected if it is in PENDING state")
	}
}

// deliverOrder changes the status of an order to DELIVERED
// An order can only be delivered if
//   - it is in APPROVED state.
//
// Once changed, the changes are irreverisble to maintain atomicity and consistency everywhere.
// In Addition to this, A Shipping Manifest is already generated for the order which cannot be redone due
// to the images and signature taken during delivery
//
// Parameters:
//   - editedOrder: Pointer to an Order object containing the edited order.
//
// Returns:
//   - error: An error object if edited status is incorrect .
func deliverOrder(editedOrder, originalOrder *models.Order) error {
	switch originalOrder.Status {
	case constants.OrderStatusApproved:
		editedOrder.Status = constants.OrderStatusDelivered
		return nil
	default:
		return errors.New("Order can only be delivered if it is in APPROVED state")
	}
}

// GetOrderWithUpdatedStatus updates the status of an order based on the edited order
//
// Parameters:
//   - editedOrder: Pointer to an Order object containing the edited order.
//   - originalOrder: Pointer to an Order object containing the original order.
//
// Returns:
//   - error: An error object if edited status is incorrect .
func GetOrderWithUpdatedStatus(editedOrder, originalOrder *models.Order) error {

	switch editedOrder.Status {
	case constants.OrderStatusApproved:
		return approveOrder(editedOrder, originalOrder)
	case constants.OrderStatusRejected:
		return rejectOrder(editedOrder, originalOrder.Status)
	case constants.OrderStatusDelivered:
		return deliverOrder(editedOrder, originalOrder)
	default:
		return nil
	}
}