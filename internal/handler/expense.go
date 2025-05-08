package handler

import (
	"encoding/json"
	"expense_tracker/internal/model"
	"expense_tracker/internal/service"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return 0, fmt.Errorf("unauthorized")
	}
	return userID, nil
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var expense model.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	creared, err := h.expenseService.CreateExpense(r.Context(), userID, expense)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(creared)
}

func (h *ExpenseHandler) GetExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expenseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid expense ID")
		return
	}

	expense, err := h.expenseService.GetExpense(r.Context(), userID, expenseID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "expense not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expenseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid expense ID")
		return
	}

	var input model.UpdateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.expenseService.UpdateExpense(r.Context(), expenseID, userID, input)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expenseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid expense ID")
		return
	}

	if err := h.expenseService.DeleteExpense(r.Context(), expenseID, userID); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) GetExpenseByPeriod(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, err := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid start time")
		return
	}
	end, err := time.Parse(time.RFC3339, r.URL.Query().Get("end"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid end time")
		return
	}

	if end.Before(start) {
		writeJSONError(w, http.StatusBadRequest, "end time must be after start time")
		return
	}

	expenses, err := h.expenseService.GetExpensesByPeriod(r.Context(), userID, start, end)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpenseHandler) GetExpenseByCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	category := r.URL.Query().Get("category")
	if category == "" {
		writeJSONError(w, http.StatusBadRequest, "category parameter is required")
		return
	}

	expenses, err := h.expenseService.GetExpensesByCategory(r.Context(), userID, category)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}
