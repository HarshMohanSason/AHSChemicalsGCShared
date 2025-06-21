package shared

import (
	"errors"
	"net/http"
	"strings"
)

// IsAuthorized verifies if the incoming HTTP request is authenticated with a valid Firebase ID token.
//
// It extracts the Bearer token from the Authorization header, verifies it using the Firebase Admin SDK,
// and returns an error if the header is malformed or token verification fails.
//
// Parameters:
//   - request: *http.Request representing the incoming HTTP request.
//
// Returns:
//   - nil if the token is valid and the request is authenticated.
//   - error if the Authorization header is missing, malformed, or if the token verification fails.
//
// Example Authorization header:
//
//	Authorization: Bearer <Firebase_ID_Token>
func IsAuthorized(request *http.Request) error {
	ctx := request.Context()

	// Extract and validate the Authorization header.
	authHeader := request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("invalid Authorization header format")
	}

	// Extract the ID token from the Authorization header.
	idToken := parts[1]
	_, err := AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return err
	}
	return nil
}

// IsAuthorizedAndAdmin verifies if the incoming HTTP request is authenticated with a valid Firebase ID token
// **and** has the custom `admin` claim set to true.
//
// It extracts the Bearer token from the Authorization header, verifies it using the Firebase Admin SDK,
// and then checks for a custom claim `"admin": true` in the token.
//
// Parameters:
//   - request: *http.Request representing the incoming HTTP request.
//
// Returns:
//   - nil if the token is valid and contains the admin claim set to true.
//   - error if:
//   - The Authorization header is missing or malformed,
//   - Token verification fails, or
//   - The token does not include the required `"admin": true` claim.
//
// Example Authorization header:
//
//	Authorization: Bearer <Firebase_ID_Token>
//
// Expected Custom Claims in Firebase Authentication:
//
//	{
//	  "admin": true
//	}
//
// Usage:
//
//	Use this function to restrict access to administrative endpoints or privileged actions.
func IsAuthorizedAndAdmin(request *http.Request) error {
	ctx := request.Context()

	// Extract and validate the Authorization header.
	authHeader := request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("invalid Authorization header format")
	}

	// Extract the ID token from the Authorization header.
	idToken := parts[1]
	token, err := AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return err
	}

	// Check if the "admin" claim exists and is set to true.
	adminClaim, ok := token.Claims["admin"].(bool)
	if !ok || !adminClaim {
		return errors.New("unauthorized: only admins are allowed to perform this operation")
	}

	return nil
}
