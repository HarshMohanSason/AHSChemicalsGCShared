package mocks

import (
	"mime/multipart"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/utils"
)

func CreateMockDeliveryInput() *models.DeliveryInput{
	return &models.DeliveryInput{
		OrderID: "AASJDKN123SAD",
		ReceivedBy: "Harsh Mohan",
		DeliveredBy: "Harsh",
		Signature: createMockDeliverySignature(),
		Images: createMockDeliveryImages(),
		DeliveredAt: time.Now(),
	}
}

func createMockDeliverySignature() multipart.File{
	file, err := utils.CreateMultipartFile("extras/mock_shipping_manifest_image.png")
	if err != nil {
		panic(err)
	}
	return file
}

func createMockDeliveryImages() []multipart.File{
	files := make([]multipart.File, 0)
	for range 2 {
		file, err := utils.CreateMultipartFile("extras/mock_shipping_manifest_image.png")
		if err != nil {
			panic(err)
		}
		files = append(files, file)
	}
	return files
}