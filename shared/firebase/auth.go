package firebase_shared

import (
	"errors"
	"net/http"
	"strings"
)

// GetUIDIfAuthorized verifies the Firebase ID token provided in the Authorization header
// of an HTTP request and returns the associated user's UID if authentication succeeds.
//
// This function is typically used to associate incoming requests with authenticated Firebase users.
//
// Parameters:
//   - request (*http.Request): The incoming HTTP request expected to contain the Authorization header.
//
// Expected Authorization header format:
//
//	Authorization: Bearer <Firebase_ID_Token>
//
// Returns:
//   - string: The Firebase user's UID if token verification is successful.
//   - error:  If the Authorization header is malformed or token verification fails.
//
// Usage Example:
//
//	uid, err := shared.GetUIDIfAuthorized(request)
//	if err != nil {
//	    // Handle unauthorized request
//	}
//	fmt.Println("User UID:", uid)
//
func GetUIDIfAuthorized(request *http.Request) (string, error) {
	ctx := request.Context()

	authHeader := request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	idToken := parts[1]

	token, err := AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}

// GetUIDIfAdmin verifies the Firebase ID token provided in the Authorization header
// of an HTTP request and returns the user's UID if the token is valid and contains
// the custom claim `"admin": true`.
//
// This function is typically used to restrict access to administrative endpoints or
// privileged operations that require elevated permissions.
//
// Parameters:
//   - request (*http.Request): The incoming HTTP request expected to contain the Authorization header.
//
// Expected Authorization header format:
//
//	Authorization: Bearer <Firebase_ID_Token>
//
// Firebase Token Requirements:
//   - The token must be a valid Firebase ID token, signed by Firebase Authentication.
//   - The token must include a custom claim `admin` with a boolean value `true`.
//
// Returns:
//   - (string): The UID (User ID) of the authenticated Firebase user.
//   - (error): Returns an error if:
//   - The Authorization header is missing, malformed, or does not follow the Bearer token format.
//   - Token verification fails via Firebase Admin SDK.
//   - The `admin` claim is missing or not set to `true`.
//
// Example:
//
//	uid, err := shared.GetUIDIfAdmin(request)
//	if err != nil {
//	    shared.WriteJSONError(response, http.StatusUnauthorized, err.Error())
//	    return
//	}
//	log.Println("Admin user UID:", uid)
//
func GetUIDIfAdmin(request *http.Request) (string, error) {
	ctx := request.Context()

	authHeader := request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	idToken := parts[1]
	token, err := AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}

	adminClaim, ok := token.Claims["admin"].(bool)
	if !ok || !adminClaim {
		return "", errors.New("unauthorized: only admins can perform this operation")
	}

	return token.UID, nil
}