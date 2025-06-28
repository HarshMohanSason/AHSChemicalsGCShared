package quickbooks

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

	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
)

// ExchangeTokenForAuthCode exchanges an OAuth 2.0 authorization code for the first time access and refresh tokens
// from the QuickBooks authorization server.
//
// It performs an HTTP POST request to the QuickBooks token endpoint with the required data and credentials.
//
// Parameters:
//   - ctx: Context for controlling request lifetime (e.g., timeouts, cancellation).
//   - receivedData: A map of required input values:
//   - "code"          → Authorization code received from QuickBooks authorization redirect
//   - "redirect_uri"  → The same redirect URI used in the authorization request
//   - "client_id"     → QuickBooks OAuth 2.0 Client ID
//   - "client_secret" → QuickBooks OAuth 2.0 Client Secret
//
// Returns:
//   - A map[string]any representing the JSON response, typically containing:
//   - "access_token"  → Token to authenticate API requests
//   - "refresh_token" → Token to obtain new access tokens
//   - "expires_in"    → Lifetime of the access token (seconds)
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
//   - "refresh_token": The refresh token obtained from the initial authorization code exchange.
//   - "client_id": QuickBooks application's Client ID.
//   - "client_secret": QuickBooks application's Client Secret.
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
func RefreshToken(ctx context.Context, refreshToken string) (map[string]any, error) {
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
//
// Returns:
//   - error: Non-nil if an error occurs during Firestore save operation
func SaveTokenToFirestore(ctx context.Context, tokenData map[string]any, authData map[string]string) error {
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

	_, err := firebase_shared.FirestoreClient.Collection("quickbooks_tokens").Doc(uid).Set(ctx, firestoreData)
	return err
}

// EnsureValidAccessToken retrieves and validates the QuickBooks access token for the specified user.
//
// If the access token is expired, it automatically refreshes the token using the stored refresh token,
// updates Firestore with the new token data, and returns the new access token.
//
// Parameters:
//   - ctx (context.Context): Context for request-scoped deadlines, cancellation, and logging.
//   - uid (string): The user ID whose QuickBooks token should be retrieved and validated.
//
// Returns:
//   - string: A valid (non-expired) QuickBooks access token.
//   - error: An error if the token does not exist, cannot be refreshed, or there are issues reading/writing Firestore.
//
// Behavior:
//  1. Retrieves the QuickBooks token document for the given user (`uid`) from Firestore.
//  2. If the document does not exist, returns an error indicating authentication is required.
//  3. If the token has expired:
//     - Attempts to refresh it using the stored `refresh_token`.
//     - Saves the new token details to Firestore.
//     - Returns the refreshed `access_token`.
//  4. If the token is still valid, returns the current `access_token`.
//
// Example Usage:
//
//	accessToken, err := shared.EnsureValidAccessToken(ctx, userID)
//	if err != nil {
//	    // Handle error, potentially requiring the user to re-authenticate
//	}
//
// Notes:
//   - Ensure that `FirestoreClient` has been properly initialized before calling this function.
//   - This function should typically be called before making authenticated requests to the QuickBooks API.
func EnsureValidAccessToken(ctx context.Context, uid string) (string, error) {
	docSnapshot, err := firebase_shared.FirestoreClient.Collection("quickbooks_tokens").Doc(uid).Get(ctx)
	if err != nil || !docSnapshot.Exists() {
		return "", fmt.Errorf("quickbooks authentication required")
	}

	docData := docSnapshot.Data()
	expiresAt, ok := docData["expires_at"].(time.Time)
	if !ok {
		return "", fmt.Errorf("invalid expires_at in token data")
	}

	//If the token is expired, generate a new one
	if time.Now().After(expiresAt) {
		refreshToken := docData["refresh_token"].(string)
		tokenResponse, err := RefreshToken(ctx, refreshToken)
		if err != nil {
			return "", err
		}
		authData := map[string]string{
			"uid":   docData["uid"].(string),
			"state": docData["state"].(string),
		}
		err = SaveTokenToFirestore(ctx, tokenResponse, authData)
		if err != nil {
			return "", err
		}
		return tokenResponse["access_token"].(string), nil
	}
	return docData["access_token"].(string), nil
}
