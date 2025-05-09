package lib

import (
	"fmt"
	"net/http"
)

type contextKey string

const UserIDkey = "userID"

func GetUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(UserIDkey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context: unauthorized")
	}
	return userID, nil
}
