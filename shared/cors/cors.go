package cors

import (
	"net/http"
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
		"http://localhost:3000":              true,
		"https://ahschemicalsdebug.web.app":  true,
		"https://azurehospitalitysupply.com": true,
	}

	// Extract the Origin header from the request.
	var origin string
	if allowedOrigins[request.Header.Get("Origin")] {
		origin = request.Header.Get("Origin")
	}

	// Set CORS headers if the origin is in the allowed list.
	if allowedOrigins[origin] {
		response.Header().Set("Access-Control-Allow-Origin", origin)
		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		response.Header().Set("Access-Control-Max-Age", "3600")
	}

	// Handle CORS preflight (OPTIONS) requests.
	if request.Method == http.MethodOptions {
		response.WriteHeader(http.StatusNoContent) // 204 No Content
		return true
	}

	return false
}