package shared

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSONSuccess(response http.ResponseWriter, statusCode int, message string, data any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	payload := map[string]any{
		"code":    statusCode,
		"message": message,
	}
	// If there is any extra data, including that data as well in the response
	if data != nil {
		payload["data"] = data
	}

	if err := json.NewEncoder(response).Encode(payload); err != nil {
		log.Printf("JSON encode response error: %v", err)
	}
}
