package services

import (
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/utils"
)

// NewDelivery creates a new delivery object from the received delivery info
//
// Parameters
//   - DeliveryInput object
//
// Returns
//   - Delivery object containing complete order details and usable images as bytes
//   - error if any
func NewDelivery(deliveryInput *models.DeliveryInput) (*models.Delivery, error) {
	err := deliveryInput.Validate()
	if err != nil {
		return nil, err
	}

	signatureBytes, err := utils.GetImageBytesFromMultiPart(deliveryInput.Signature)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the signature file, please try again: %s", err.Error())
	}
	if signatureBytes == nil {
		return nil, fmt.Errorf("No signature found")
	}

	deliveryImagesBytes := make([][]byte, 0)
	for _, deliveryImage := range deliveryInput.Images {
		imageBytes, err := utils.GetImageBytesFromMultiPart(deliveryImage)
		if err != nil {
			return nil, fmt.Errorf("Error while reading the image file, please try again: %s", err.Error())
		}
		if imageBytes == nil {
			return nil, fmt.Errorf("No image found")
		}
		deliveryImagesBytes = append(deliveryImagesBytes, imageBytes)
	}

	return &models.Delivery{
		DeliveredBy:    deliveryInput.DeliveredBy,
		ReceivedBy:     deliveryInput.ReceivedBy,
		Signature:      signatureBytes,
		DeliveryImages: deliveryImagesBytes,
		DeliveredAt:    deliveryInput.DeliveredAt,
	}, nil
}
