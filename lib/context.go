package lib

import (
	"fmt"
	"net/http"
)

type contextKey string

// UserIDkey is the key used to store and retrieve the user ID from the request context.
const UserIDkey = "userID"

// GetUserIDFromContext extracts the user ID from the request context.
// It returns the user ID as an integer if present, or an error if not found.
//
// Parameters:
// - r: the HTTP request containing the context.
//
// Returns:
// - int: the user ID stored in the context.
// - error: an error indicating that the user ID was not found (unauthorized).
//
// Usage:
// The user ID must be previously set in the context using the same key (UserIDkey).
func GetUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(UserIDkey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context: unauthorized")
	}
	return userID, nil
}
