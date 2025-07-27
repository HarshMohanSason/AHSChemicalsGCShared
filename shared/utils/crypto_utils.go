//utils package contains common utility functions that are used across multiple packages in the application
package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomSecret generates a cryptographically secure random string,
// suitable for use as OAuth 'state' parameters, API secrets, etc.
//
// Returns:
//   - string: A URL-safe base64 encoded random secret of 32 bytes (typically ~43 characters).
//
// Example:
//   secret := generateRandomSecret()
func GenerateRandomSecret() (string, error) {
	b := make([]byte, 32) 
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}