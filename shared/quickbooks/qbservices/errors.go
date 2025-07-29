package qbservices

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
)

func ReturnErrorFromQBResp(respBody []byte, apiName string) error {
	var errResp qbmodels.QBErrorResponse
	err := json.Unmarshal(respBody, &errResp)
	if err != nil {
		return errors.New("Error occurred parsing the quickbooks error response body of the api: " + apiName)
	}
	return fmt.Errorf("Error receieved in response while making a request to %s with Code: %s, Message: %s", apiName, errResp.Fault.Error[0].Code, errResp.Fault.Error[0].Message)
}
