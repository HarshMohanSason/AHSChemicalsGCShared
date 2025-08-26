package tests

import (
	"context"
	"testing"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/repositories"
)

func TestContactUsSaveToFirestore(t *testing.T) {
	ip := "123.12.21.0"
	ctx := context.Background()

	contactForm := &models.ContactUsForm{
		Name:      "test",
		Email:     "test",
		Phone:     "2312123123",
		Location:  "2040 N Preisker lane",
		Message:   "test",
		Timestamp: time.Now().UTC().Add(48 * time.Hour),
	}

	err := repositories.SaveContactUsToFirestore(contactForm, ip, ctx)
	if err != nil {
		t.Error(err)
	}
}
