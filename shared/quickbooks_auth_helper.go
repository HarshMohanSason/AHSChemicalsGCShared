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
)

// ExchangeTokenForAuthCode exchanges an OAuth 2.0 authorization code for the first time access and refresh tokens
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
func ExchangeTokenForAuthCode(ctx context.Context, authCode string) (map[string]any, error) {
	tokenURL := "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"

	//Set the url values 
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", QUICKBOOKS_AUTH_CALLBACK_URL)

	authStr := fmt.Sprintf("%s:%s", QUICKBOOKS_CLIENT_ID, QUICKBOOKS_CLIENT_SECRET) 
	baseAuth := base64.StdEncoding.EncodeToString([]byte(authStr))

	//Create a new request with encoded url as the body
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

// RefreshToken exchanges a valid QuickBooks refresh token for a new access token.
//
// This function implements the OAuth 2.0 "refresh_token" grant type, allowing the application
// to obtain a new access token when the current one expires.
//
// Parameters:
//   - ctx (context.Context): Context for controlling request lifetime (timeouts, cancellations, etc.).
//   - refreshTokenData (map[string]string): Map containing:
//       - "refresh_token": The refresh token obtained from the initial authorization code exchange.
//       - "client_id": QuickBooks application's Client ID.
//       - "client_secret": QuickBooks application's Client Secret.
//
// Returns:
//   - (map[string]any): The response from QuickBooks containing the new access token and related metadata.
//   - (error): Any error encountered during the HTTP request or JSON parsing.
//
// Example usage:
//
//	tokenData, err := RefreshToken(ctx, map[string]string{
//	    "refresh_token": storedRefreshToken,
//	    "client_id":     QUICKBOOKS_CLIENT_ID,
//	    "client_secret": QUICKBOOKS_CLIENT_SECRET,
//	})
//
// Response JSON Example:
//
//	{
//	  "access_token": "AAABBBCCC...",
//	  "refresh_token": "DDDDEEEEFFFF...",
//	  "expires_in": 3600,
//	  "token_type": "bearer",
//	  "x_refresh_token_expires_in": 8726400
//	}
func RefreshToken(ctx context.Context, refreshToken string) (map[string]any, error){
	tokenURL := "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	authStr := fmt.Sprintf("%s:%s", QUICKBOOKS_CLIENT_ID, QUICKBOOKS_CLIENT_SECRET)
	baseAuth := base64.StdEncoding.EncodeToString([]byte(authStr))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", baseAuth))
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
		return nil, fmt.Errorf("refresh token failed: %s", body)
	}

	var tokenResp map[string]any
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return tokenResp, nil
}
// SaveTokenToFirestore saves the provided token data to Firestore with calculated expiration timestamps.
//
// Parameters:
//   - ctx: Context for request lifetime
//   - tokenData: Token response data from QuickBooks (access_token, refresh_token, etc.)
//   - authData: Metadata map including "uid" and "state"
//   - firestoreClient: Firestore client instance
//
// Returns:
//   - error: Non-nil if an error occurs during Firestore save operation
func SaveTokenToFirestore(ctx context.Context, tokenData map[string]any, authData map[string]string) error{
	uid, ok := authData["uid"]
	if !ok || uid == "" {
		return fmt.Errorf("missing UID in authData")
	}
	// Calculate expiration timestamp
	obtainedAt := time.Now()
	expiresInSec := int64(tokenData["expires_in"].(float64))
	expiresAt := obtainedAt.Add(time.Duration(expiresInSec) * time.Second)

	firestoreData := map[string]any{
		"access_token":  tokenData["access_token"],
		"refresh_token": tokenData["refresh_token"],
		"expires_in":    expiresInSec,
		"obtained_at":   obtainedAt,
		"expires_at":    expiresAt,
		"token_type":    tokenData["token_type"],
		"scope":         tokenData["scope"],
		"state":         authData["state"],
	}

	_ , err := FirestoreClient.Collection("quickbooks_tokens").Doc(uid).Set(ctx, firestoreData)
	return err
}
