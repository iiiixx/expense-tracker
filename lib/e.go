package lib

import (
	"encoding/json"
	"net/http"
)

// WriteJSONError sends a JSON-formatted error response to the client.
// It sets the HTTP status code and returns a JSON object with an "error" field containing the message.
//
// Parameters:
// - w: the http.ResponseWriter to write the response.
// - status: the HTTP status code to set in the response.
// - message: the error message to include in the JSON response.
//
// Example JSON response:
//
//	{
//	  "error": "description of the error"
//	}
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
