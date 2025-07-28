package repositories

import (
	"context"
	"errors"
	"time"

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SaveContactUsToFirestore saves contact us form to firestore
//
// Parameters:
// c *models.ContactUsForm - contact us form pointer to the struct
// ip string - ip address of the client
// ctx context.Context - context of the request
//
// Returns:
// error - error if any
func SaveContactUsToFirestore(c *models.ContactUsForm, ip string, ctx context.Context) error {

	docSnapshot, err := firebase_shared.FirestoreClient.Collection("contact_us").Doc(ip).Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			_, err = firebase_shared.FirestoreClient.Collection("contact_us").Doc(ip).Set(ctx, c)
			return err
		}
		return err
	}

	if !docSnapshot.Exists() {
		_, err = firebase_shared.FirestoreClient.Collection("contact_us").Doc(ip).Set(ctx, c)
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

	_, err = firebase_shared.FirestoreClient.Collection("contact_us").Doc(ip).Set(ctx, c)
	if err != nil {
		return err
	}
	return nil
}