package shared

import (
	"encoding/json"
	"net/http"
)

type FirebaseErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

func WriteJSONError(response http.ResponseWriter, statusCode int, message string){
	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(statusCode)
	_ = json.NewEncoder(response).Encode(map[string]any{"code": statusCode, "message": message})
}