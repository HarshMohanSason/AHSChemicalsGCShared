package repositories

import (
	"context"
	"errors"
	"time"

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
)

func SaveContactUsToFirestore(c *models.ContactUsForm, ip string, ctx context.Context) error {

	docSnapshot, err := firebase_shared.FirestoreClient.Collection("contact_us").Doc(ip).Get(ctx)
	if err != nil {
		return err
	}

	var oldContactUs models.ContactUsForm
	err = docSnapshot.DataTo(&oldContactUs)
	if err != nil {
		return err
	}

	if time.Since(oldContactUs.Timestamp) < time.Hour*24 {
		return errors.New("You need to wait 24 hours before submitting another contact us request")
	}

	//Save the contact us form
	_, err = firebase_shared.FirestoreClient.Collection("contact_us").Doc(ip).Set(ctx, c)
	if err != nil {
		return err
	}
	return nil
}