package qbservices

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

	"cloud.google.com/go/firestore"
	firebase_shared "github.com/HarshMohanSason/AHSChemicalsGCShared/shared/firebase"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks/qbmodels"
)

// ExchangeTokenForAuthCode exchanges an authorization code fetched from the auth callback url for the token reponse from QuickBooks.
func ExchangeTokenForAuthCode(ctx context.Context, authCode string) (*qbmodels.QBReponseToken, error) {
	tokenURL := "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", quickbooks.QUICKBOOKS_AUTH_CALLBACK_URL)

	authStr := fmt.Sprintf("%s:%s", quickbooks.QUICKBOOKS_CLIENT_ID, quickbooks.QUICKBOOKS_CLIENT_SECRET)
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
		return nil, ReturnQBOAuthError(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", body)
	}

	var tokenResp qbmodels.QBReponseToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// refreshToken uses a valid refresh token to obtain a new access token.
func refreshToken(ctx context.Context, refreshToken string) (*qbmodels.QBReponseToken, error) {
	tokenURL := "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	authStr := fmt.Sprintf("%s:%s", quickbooks.QUICKBOOKS_CLIENT_ID, quickbooks.QUICKBOOKS_CLIENT_SECRET)
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
		return nil, ReturnQBOAuthError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh token failed: %s", body)
	}

	var tokenResp qbmodels.QBReponseToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

// EnsureValidAccessToken retrieves and validates the QuickBooks access token and related OAuth data
// for the specified user from Firestore.
func EnsureValidAccessToken(ctx context.Context, uid string) (*qbmodels.QBReponseToken, error) {
	docSnapshot, err := firebase_shared.FirestoreClient.Collection("quickbooks_tokens").Doc(uid).Get(ctx)
	if err != nil || !docSnapshot.Exists() {
		return nil, fmt.Errorf("quickbooks authentication required")
	}

	var originalToken qbmodels.QBReponseToken
	err = docSnapshot.DataTo(&originalToken)
	if err != nil {
		return nil, err
	}

	if originalToken.IsExpired() {
		newToken, err := refreshToken(ctx, originalToken.RefreshToken)
		if err != nil {
			return nil, err
		}
		err = SaveTokenToFirestore(ctx, newToken, uid)
		if err != nil {
			return nil, err
		}
		return newToken, nil
	} else {
		return &originalToken, nil
	}
}

// SaveTokenToFirestore saves the QBTokenResponse object to firestore with timestamps.
func SaveTokenToFirestore(ctx context.Context, t *qbmodels.QBReponseToken, uid string) error {
	t.SetObtainedAt()
	t.SetExpiresAt()

	_, err := firebase_shared.FirestoreClient.Collection("quickbooks_tokens").Doc(uid).Set(ctx, t, firestore.MergeAll)
	return err
}