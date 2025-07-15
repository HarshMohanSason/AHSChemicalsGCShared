package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSONSuccess sends a standardized JSON success response to the client.
//
// This helper function sets the appropriate Content-Type header, writes the specified HTTP status code,
// and encodes a JSON object containing the status code, message, and optional data payload.
//
// Parameters:
//   - response: http.ResponseWriter to write the response to.
//   - statusCode: HTTP status code to send (e.g., 200, 201, etc.).
//   - message: Human-readable success message.
//   - data: Optional payload to include in the response body (can be any Go type, or nil).
//
// Example JSON Response:
//   {
//     "code": 200,
//     "message": "Operation successful",
//     "data": {...} // optional
//   }
//
// Logs:
//   - Logs any errors encountered during JSON encoding of the response.
func WriteJSONSuccess(response http.ResponseWriter, statusCode int, message string, data any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	payload := map[string]any{
		"code":    statusCode,
		"message": message,
	}

	// Include optional data in the response if provided.
	if data != nil {
		payload["data"] = data
	}

	if err := json.NewEncoder(response).Encode(payload); err != nil {
		log.Printf("JSON encode response error: %v", err)
	}
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