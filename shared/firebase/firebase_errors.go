package firebase_shared

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// FirebaseErrorResponse represents the structure of an error response returned by Firebase Admin SDK.
//
// This struct is used specifically to parse and extract meaningful error details
// from JSON error responses originating from Firebase services.
type FirebaseErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`    // HTTP status code of the error.
		Message string `json:"message"` // Summary message of the error.
		Errors  []struct {
			Message string `json:"message"` // Detailed error message.
			Domain  string `json:"domain"`  // The service-specific domain of the error.
			Reason  string `json:"reason"`  // Specific reason for the error.
		} `json:"errors"`
	} `json:"error"`
}

// ExtractFirebaseErrorFromResponse attempts to extract a FirebaseErrorResponse from a given error.
//
// Firebase sometimes returns error responses as JSON strings prefixed with additional data.
// This function locates the JSON portion, unmarshals it into a FirebaseErrorResponse struct, and returns it.
//
// Parameters:
//   - err: The original error object returned by Firebase.
//
// Returns:
//   - A pointer to FirebaseErrorResponse if extraction is successful.
//   - nil if the JSON portion is not found.
//   - If JSON is malformed, returns a FirebaseErrorResponse with empty or partial fields.
func ExtractFirebaseErrorFromResponse(err error) *FirebaseErrorResponse {
	errString := err.Error()

	// Locate the start of the JSON object within the error string.
	start := strings.Index(errString, "{")
	var firebaseError FirebaseErrorResponse

	// If no JSON object is found, return nil.
	if start == -1 {
		return nil
	}

	// Extract and attempt to unmarshal the JSON portion of the error.
	jsonPart := errString[start:]
	if unmarshalErr := json.Unmarshal([]byte(jsonPart), &firebaseError); unmarshalErr != nil {
		return &firebaseError // Return partial object even if unmarshaling fails.
	}
	return &firebaseError
}

// WriteJSONError sends a standardized JSON error response to the client.
//
// This helper function sets the appropriate Content-Type, writes the HTTP status code,
// and encodes a JSON object containing the status code and error message.
//
// Parameters:
//   - response: The http.ResponseWriter to write the response to.
//   - statusCode: The HTTP status code to send (e.g., 400, 404, 500).
//   - message: A descriptive error message for the client.
//
// Logs:
//   - Logs any errors that occur while writing the response.
func WriteJSONError(response http.ResponseWriter, statusCode int, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	if err := json.NewEncoder(response).Encode(map[string]any{"code": statusCode, "message": message}); err != nil {
		log.Printf("Error writing the error response: %v", err)
	}
}