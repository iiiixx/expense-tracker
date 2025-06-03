package handler

import (
	"encoding/json"
	"expense_tracker/internal/model"
	"expense_tracker/internal/service"
	"expense_tracker/lib"
	"net/http"
)

// UserHandler handles HTTP requests related to user operations.
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler with the given UserService.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UpdateUsername handles the HTTP request to update a user's username.
// Possible HTTP responses:
// - 200 OK: Username updated successfully.
// - 400 Bad Request: Invalid request body or update error.
// - 401 Unauthorized: User authentication failed.
func (h *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var input model.UpdateUsernameInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedUser, err := h.userService.UpdateUsername(r.Context(), userID, &input)
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser handles the HTTP request to delete the authenticated user.
// Possible HTTP responses:
// - 204 No Content: User deleted successfully.
// - 400 Bad Request: Deletion error.
// - 401 Unauthorized: User authentication failed.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := h.userService.DeleteUser(r.Context(), userID); err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetProfile handles the HTTP request to retrieve the authenticated user's profile.
// Possible HTTP responses:
// - 200 OK: Profile retrieved successfully.
// - 401 Unauthorized: User authentication failed.
// - 500 Internal Server Error: Failed to retrieve profile.
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		lib.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
}
