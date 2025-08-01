package qbservices

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ReturnErrorFromQBResp(body, "ExchangeTokenForAuthCode")
	}

	var tokenResp qbmodels.QBReponseToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// refreshToken uses a valid refresh token to obtain a new access token. If the refresh token has expired, it will return an error.
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
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ReturnErrorFromQBResp(body, "RefreshToken")
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
		//Check if refresh token is expired or not. If yes, return an error since user needs to re-login to quickbooks.
		if originalToken.IsRefreshTokenExpired(){
			return nil, fmt.Errorf("Quickbooks session has expired. Please login again to get a new access token.")
		}
		newToken, err := refreshToken(ctx, originalToken.RefreshToken)
		if err != nil {
			return nil, err
		}
		//newToken does not return the realmId and state.
		newToken.SetRealmID(originalToken.RealmId)
		newToken.SetState(originalToken.State)
		
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

	_, err := firebase_shared.FirestoreClient.Collection("quickbooks_tokens").Doc(uid).Set(ctx, t.ToMap(), firestore.MergeAll)
	return err
}

//Only used in webhooks since there is no way to pass the user uid when making changes in quickbooks.
//So this fetches the uid from the first document found in the quickbooks_tokens.
func GetTokenUIDFromFirestore(ctx context.Context) (string, error) {
	docRefs, err := firebase_shared.FirestoreClient.Collection("quickbooks_tokens").Documents(ctx).GetAll()
	if err != nil{
		return "", err
	}
	if len(docRefs) == 0 {
		return "", errors.New("No admin token for quickbooks found. Please re authenticate quickbooks")
	}
	return docRefs[0].Ref.ID, nil
}