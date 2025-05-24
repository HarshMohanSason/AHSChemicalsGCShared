package shared

import (
	"net/http"
)

func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		
		corsEnabledFunction(response, request)
		
		next.ServeHTTP(response, request)
	})
}

func corsEnabledFunction(response http.ResponseWriter, request *http.Request) {
	
	//Pre defined allowed origins
	allowedOrigins := map[string]bool{
		"http://localhost:3000": true,
		"https://ahschemicalsdebug.web.app": true,
	}

	//Check if the origin exists
	var origin string
	if allowedOrigins[request.Header.Get("Origin")] {
		origin = request.Header.Get("Origin")
	}

	// Set CORS headers for the preflight request
	if request.Method == http.MethodOptions {
		response.Header().Set("Access-Control-Allow-Origin", origin)
		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		response.Header().Set("Access-Control-Max-Age", "3600")
		response.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	response.Header().Set("Access-Control-Allow-Origin", origin)
	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
}