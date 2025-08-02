package tests

import (
	"context"
	"testing"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/gcp"
)

func TestPublishMessage(t *testing.T) {

	err := gcp.PublishMessage(context.Background(), []byte("Harsh Mohan Sason"))
	if err != nil {
		t.Error(err)
	}
	t.Log("Message Published")
}
