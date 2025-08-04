package models

import (
	"errors"
	"mime/multipart"
	"time"
)

type DeliveryInput struct {
	OrderID     string           
	ReceivedBy  string          
	DeliveredBy string          
	Signature   multipart.File  
	Images      []multipart.File      
}

func (d *DeliveryInput) SetOrderID(orderID string) {
	d.OrderID = orderID
}
func (d *DeliveryInput) SetReceivedBy(receivedBy string) {
	d.ReceivedBy = receivedBy
}
func (d *DeliveryInput) SetDeliveredBy(deliveredBy string) {
	d.DeliveredBy = deliveredBy
}
func (d *DeliveryInput) SetSignature(signature multipart.File) {
	d.Signature = signature
}
func (d *DeliveryInput) SetImages(images []multipart.File) {
	d.Images = images
}

func (d *DeliveryInput) Validate() error {
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
	return nil
}

type Delivery struct {
	Order          *Order
	ReceivedBy     string
	DeliveredBy    string
	Signature      []byte
	DeliveryImages [][]byte
	DeliveredAt    time.Time
}