package handler

import (
	"encoding/json"
	"expense_tracker/internal/model"
	"expense_tracker/internal/service"
	"expense_tracker/lib"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

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
