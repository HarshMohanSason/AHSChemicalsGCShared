package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
)

func TestIsAuthroized(t *testing.T){
	
	jwtToken := "enterthejwthere:)"
	bearer := fmt.Sprintf("Bearer %s", jwtToken)
	
	//Stimulate a new request 
	req := httptest.NewRequest(http.MethodPost, "https://us-west2-ahschemicalsprod.cloudfunctions.net/create-account", nil)
	req.Header.Set("Authorization", bearer)
    req.Header.Set("Content-Type", "application/json")

    //Error assertion
    err := shared.IsAuthorized(req)
    if err != nil{
    	t.Errorf("An error occured %v", err)
    }
}