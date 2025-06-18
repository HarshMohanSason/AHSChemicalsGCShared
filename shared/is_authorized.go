package shared

import (
	"errors"
	"net/http"
	"strings"
)

// Checks custom claims to detect if an admin or not
func IsAuthorized(request *http.Request) error {

	ctx := request.Context()

	//Process the authroization header
	authHeader := request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("Invalid auth header")
	}

	//Get the idToken
	idToken := parts[1]
	token, err := AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return err
	}

	//Check if the user is admin or not
	if !token.Claims["admin"].(bool) {
		return errors.New("Unauthorized. Only admins are allowed to perform this operation")
	}

	return nil
}
