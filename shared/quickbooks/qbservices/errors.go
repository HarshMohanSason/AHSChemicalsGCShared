package qbservices

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
)

// Error code to know if the quickbooks session was expired or not. Quickbooks does not allow to reauthenticate programtiically. 
// User needs to manually re-login into quickbooks via email and password.
// Currently used in `quickbooks-webhook-entity-processor` cloud event.
var ErrQuickBooksSessionExpired = errors.New("quickbooks session has expired, user must re-login")

func ReturnErrorFromQBResp(respBody []byte, apiName string) error {
	var errResp qbmodels.QBErrorResponse
	err := json.Unmarshal(respBody, &errResp)
	if err != nil {
		return errors.New("Error occurred parsing the quickbooks error response body of the api: " + apiName)
	}
	return fmt.Errorf("Error receieved in response while making a request to %s with Code: %s, Message: %s", apiName, errResp.Fault.Error[0].Code, errResp.Fault.Error[0].Message)
}