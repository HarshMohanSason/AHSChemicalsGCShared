package shared

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
)

// ExchangeTokenForAuthCode exchanges an OAuth 2.0 authorization code for access and refresh tokens
// from the QuickBooks authorization server.
//
// It performs an HTTP POST request to the QuickBooks token endpoint with the required data and credentials.
//
// Parameters:
//   - ctx: Context for controlling request lifetime (e.g., timeouts, cancellation).
//   - receivedData: A map of required input values:
//       • "code"          → Authorization code received from QuickBooks authorization redirect
//       • "redirect_uri"  → The same redirect URI used in the authorization request
//       • "client_id"     → QuickBooks OAuth 2.0 Client ID
//       • "client_secret" → QuickBooks OAuth 2.0 Client Secret
//
// Returns:
//   - A map[string]any representing the JSON response, typically containing:
//       • "access_token"  → Token to authenticate API requests
//       • "refresh_token" → Token to obtain new access tokens
//       • "expires_in"    → Lifetime of the access token (seconds)
//   - An error, if any occurred during the request, response handling, or decoding.
//
// Example usage:
//
//	receivedData := map[string]string{
//		"code":          "<authorization_code>",
//		"redirect_uri":  "<your_redirect_uri>",
//		"client_id":     "<your_client_id>",
//		"client_secret": "<your_client_secret>",
//	}
//	tokens, err := shared.ExchangeTokenForAuthCode(ctx, receivedData)
//
// Notes:
//   - Requires HTTPS for secure transmission.
//   - If the response has a non-200 HTTP status code, the error will include the response body for debugging.
//
// OAuth 2.0 Reference:
//   - https://developer.intuit.com/app/developer/qbo/docs/develop/authentication-and-authorization
func ExchangeTokenForAuthCode(ctx context.Context, receivedData map[string]string) (map[string]any, error) {
	tokenURL := "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", receivedData["code"])
	data.Set("redirect_uri", receivedData["redirect_uri"])

	authStr := fmt.Sprintf("%s:%s", receivedData["client_id"], receivedData["client_secret"])
	baseAuth := base64.StdEncoding.EncodeToString([]byte(authStr))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	authValue := fmt.Sprintf("Basic %s", baseAuth)
	req.Header.Set("Authorization", authValue)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", body)
	}

	var tokenResp map[string]any
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return tokenResp, nil
}