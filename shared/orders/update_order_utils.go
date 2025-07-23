package orders

import (
	"context"
	"errors"
	"time"
)

//approveOrder changes the status of an order to APPROVED
//An order can only be approved if 
//  - it is in PENDING state
//  - it is in REJECTED state and no more than 30 days have passed since rejection
//
//Parameters:
//  - editedOrder: Pointer to an Order object containing the edited order.
//  - originalOrder: Pointer to an Order object containing the original order.
//
//Returns: 
//  - error: An error object if edited status is incorrect .
func approveOrder(editedOrder *Order, originalOrder *Order) error {
	switch originalOrder.Status {
	case OrderStatusPending:
		editedOrder.Status = OrderStatusApproved
		return nil
	case OrderStatusRejected:
		if time.Since(originalOrder.UpdatedAt) > 30*24*time.Hour {
			return errors.New("Cannot approve: more than 30 days have passed since rejection. Place a new order instead")
		}
		editedOrder.Status = OrderStatusApproved
		return nil
	default:
		return errors.New("Order can only be approved if it is in PENDING or REJECTED state")
	}
}

//rejectOrder changes the status of an order to REJECTED.
//An order can only be rejected if
//  - it is in PENDING state
//
//Parameters:
//  - editedOrder: Pointer to an Order object containing the edited order.
//  - originalOrder: Pointer to an Order object containing the original order.
//
//Returns:
//  - error: An error object if edited status is incorrect .
func rejectOrder(editedOrder *Order, originalOrderStatus string) error {
	switch originalOrderStatus {
		case OrderStatusPending: 
			editedOrder.Status = OrderStatusRejected
			return nil
		default:
			return errors.New("Order can only be rejected if it is in PENDING state")
	}
}

//deliverOrder changes the status of an order to DELIVERED
//An order can only be delivered if
//  - it is in APPROVED state.
//Once changed, the changes are irreverisble to maintain atomicity and consistency everywhere. 
//In Addition to this, A Shipping Manifest is already generated for the order which cannot be redone due 
//to the images and signature taken during delivery 
//
//Parameters:
//  - editedOrder: Pointer to an Order object containing the edited order.
//
//Returns: 
//  - error: An error object if edited status is incorrect .
func deliverOrder(editedOrder *Order, originalOrder *Order) error {
	switch originalOrder.Status {
		case OrderStatusApproved:
			editedOrder.Status = OrderStatusDelivered
			return nil
		default:
			return errors.New("Order can only be delivered if it is in APPROVED state")
	}
}

//updateOrderStatus updates the status of an order based on the edited order
//
//Parameters:
//  - editedOrder: Pointer to an Order object containing the edited order.
//  - originalOrder: Pointer to an Order object containing the original order.
//
//Returns:
//  - error: An error object if edited status is incorrect .
func updateOrderStatus(editedOrder *Order, originalOrder *Order) error {
	
	switch(editedOrder.Status) {
		case OrderStatusApproved:
			return approveOrder(editedOrder, originalOrder)
		case OrderStatusRejected:
			return rejectOrder(editedOrder, originalOrder.Status)
		case OrderStatusDelivered:
			return deliverOrder(editedOrder, originalOrder)
		default: 
			return nil
	}
}

//checkIfOrderItemsChanged checks if the items in the order have any changes.
//
// Parameters:
//   - editedOrder: Pointer to an Order object containing the edited order.
//   - originalOrder: Pointer to an Order object containing the original order.
//
// Returns:
//   - bool: A boolean value indicating whether the items in the order have any changes.
func checkIfOrderItemsChanged(editedOrder *Order, originalOrder *Order) bool {

	for i, item := range editedOrder.Items {
		if item.UnitPrice != originalOrder.Items[i].UnitPrice {
			return true
		}
		if item.Quantity != originalOrder.Items[i].Quantity {
			return true
		}
	}
	return false
}

//GetUpdatableOrder checks if the order has any changes made from the original order before updating.
//If the order does not any changes, an error is returned, else the modified order is returned
//
// Parameters:
//   - editedOrder: Pointer to an Order object containing the edited order.
//   - ctx: A context object used for asynchronous operations.
//
// Returns
//   - error: An error object if edited status is incorrect .
func GetUpdatableOrder(editedOrder *Order, ctx context.Context) error {

	originalOrder, err := FetchOrderFromFirestore(editedOrder.ID, ctx)
	if err != nil {
		return err
	}
	if originalOrder.Status == editedOrder.Status || !checkIfOrderItemsChanged(editedOrder, originalOrder) {
		return errors.New("No new changes were detected. Please make changes before updating the order")
	}
	//Update the order status
	updateOrderStatus(editedOrder, originalOrder)
	return nil
}