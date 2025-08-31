package models

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"mime/multipart"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/utils"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
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

func (d *Delivery) GetDeliveredAtLocalTime() time.Time {
	localTime, err := utils.ConvertUTCToLocalTimeZoneWithFormat(d.DeliveredAt, d.Order.TimeZone)
	if err != nil {
		return d.DeliveredAt
	}
	return localTime
}

func (d *Delivery) GetCorrectlyRotatedImages() [][]byte {
	deliveryImages := make([][]byte, 0)
	for _, imageBytes := range d.DeliveryImages {
		img, _, err := image.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			continue
		}

		//Get orientation of the image
		x, err := exif.Decode(bytes.NewReader(imageBytes))
		orientation := 1 // default
		if err == nil {
			tag, err := x.Get(exif.Orientation)
			if err == nil {
				o, err := tag.Int(0)
				if err == nil {
					orientation = o
				}
			}
		}

		img = utils.FixImageOrientation(img, orientation)
		img = imaging.Clone(img) // Force 8-bit NRGBA
		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			continue
		}
		deliveryImages = append(deliveryImages, buf.Bytes())
	}
	if len(deliveryImages) == 0 {
		return d.DeliveryImages
	}
	return deliveryImages
}
