package utils

import (
	"net/http"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/constants"
)

// CorsEnabledFunction handles Cross-Origin Resource Sharing (CORS) for HTTP requests.
//
// This function verifies if the request's Origin header is in the list of allowed origins.
// If allowed, it sets appropriate CORS headers on the HTTP response to enable cross-origin access.
//
// Additionally, it handles preflight (OPTIONS) requests by returning HTTP 204 No Content.
//
// Parameters:
//   - response: http.ResponseWriter used to write HTTP headers and responses.
//   - request: *http.Request representing the incoming HTTP request.
//
// Returns:
//   - true: If the request is an OPTIONS preflight request and has been handled.
//   - false: If the request is not an OPTIONS preflight request (normal processing should continue).
func CorsEnabledFunction(response http.ResponseWriter, request *http.Request) bool {
	allowedOrigins := map[string]bool{
		constants.CorsAllowOriginDebug:      true,
		constants.CorsAllowOriginStaging:    true,
		constants.CorsAllowOriginProduction: true,
	}

	var origin string
	if allowedOrigins[request.Header.Get("Origin")] {
		origin = request.Header.Get("Origin")
	}

	if allowedOrigins[origin] {
		response.Header().Set("Access-Control-Allow-Origin", origin)
		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		response.Header().Set("Access-Control-Max-Age", "3600")
	}

	if request.Method == http.MethodOptions {
		response.WriteHeader(http.StatusNoContent) // 204 No Content
		return true
	}

	return false
}