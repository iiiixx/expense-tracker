package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func GetUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return 0, fmt.Errorf("unauthorized")
	}
	return userID, nil
}
