package models

import (
	"errors"
	"mime/multipart"
	"time"
)

type Delivery struct {
	OrderID     string           `json:"orderId"`
	ReceivedBy  string           `json:"receivedBy"`  //Name of the person who received the order
	DeliveredBy string           `json:"deliveredBy"` //Name of the person who delivered the order
	Signature   multipart.File   `json:"signature"`   //Signature of the person who received the order
	Images      []multipart.File `json:"images"`      //Images of the product when delivered
	DeliveredAt time.Time        `json:"deliveredAt"`
}

func (d *Delivery) Validate() error {
	if d.OrderID == "" {
		return errors.New("No order id was found when saving delivery. Please retry submission again")
	}
	if d.ReceivedBy == "" {
		return errors.New("No name was found when saving delivery. Please retry submission again")
	}
	if d.DeliveredBy == "" {
		return errors.New("No name was found when saving delivery. Please retry submission again")
	}
	if d.Signature == nil {
		return errors.New("No signature was found when saving delivery. Please retry submission again")
	}
	if len(d.Images) == 0 {
		return errors.New("No images were found when saving delivery. At least one image required. Please retry submission again")
	}
	if d.DeliveredAt.IsZero() {
		return errors.New("No date was found when saving delivery. Please retry submission again")
	}
	return nil
}