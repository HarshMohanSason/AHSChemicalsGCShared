package qbservices

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/utils"
)

func GetOrderQBEstimate(order *models.Order) (estimate *qbmodels.QBEstimate, err error) {
	reqURL := quickbooks.QUICKBOOKS_GET_ESTIMATE_URL

	newEstimate := qbmodels.NewQBEstimate(order)
	
	//Convert to JSON bytes
	estimateJSON, err := json.Marshal(newEstimate)
	if err != nil {
		return nil, err
	}

	//Create request
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(estimateJSON))
	if err != nil {
		return nil, err
	}

	//Do request
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to get item from quickbooks" + string(respBody))
	}

	var respPaylod utils.SuccessPaylod[string]
	err = json.Unmarshal(respBody, &respPaylod)
	if err != nil {
		return nil, errors.New("Could not get item from quickbooks" + string(respBody))
	}

	var quickbooksItem qbmodels.QBEstimateResponse
	err = json.Unmarshal([]byte(respPaylod.Data), &quickbooksItem)
	if err != nil {
		return nil, err
	}

	return quickbooksItem.Estimate, nil
}
