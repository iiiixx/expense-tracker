package handler

import (
	"encoding/json"
	"expense_tracker/internal/model"
	"expense_tracker/internal/service"
	"expense_tracker/lib"
	"net/http"
	"strconv"
	"time"
)

// ExpenseHandler handles HTTP requests related to expense operations.
type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

// NewExpenseHandler creates a new ExpenseHandler with the given ExpenseService.
func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// CreateExpense handles the HTTP request to create a new expense for the authenticated user.
// Possible HTTP responses:
// - 201 Created: Expense created successfully.
// - 400 Bad Request: Invalid request body or creation error.
// - 401 Unauthorized: User authentication failed.
func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var expense model.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	creared, err := h.expenseService.CreateExpense(r.Context(), userID, expense)
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(creared)
}

// GetExpense handles the HTTP request to retrieve a specific expense by its ID for the authenticated user.
// Possible HTTP responses:
// - 200 OK: Expense retrieved successfully.
// - 400 Bad Request: Invalid expense ID.
// - 401 Unauthorized: User authentication failed.
// - 404 Not Found: Expense not found.
func (h *ExpenseHandler) GetExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expenseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid expense ID")
		return
	}

	expense, err := h.expenseService.GetExpense(r.Context(), userID, expenseID)
	if err != nil {
		lib.WriteJSONError(w, http.StatusNotFound, "expense not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

// UpdateExpense handles the HTTP request to update an existing expense by its ID for the authenticated user.
// Possible HTTP responses:
// - 200 OK: Expense updated successfully.
// - 400 Bad Request: Invalid expense ID, request body, or update error.
// - 401 Unauthorized: User authentication failed.
func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expenseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid expense ID")
		return
	}

	var input model.UpdateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.expenseService.UpdateExpense(r.Context(), expenseID, userID, &input)
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteExpense handles the HTTP request to delete an expense by its ID for the authenticated user.
// Possible HTTP responses:
// - 204 No Content: Expense deleted successfully.
// - 400 Bad Request: Invalid expense ID or deletion error.
// - 401 Unauthorized: User authentication failed.
func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expenseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid expense ID")
		return
	}

	if err := h.expenseService.DeleteExpense(r.Context(), expenseID, userID); err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetExpensesList handles the HTTP request to retrieve a list of all expenses for the authenticated user.
// Possible HTTP responses:
// - 200 OK: Expenses list retrieved successfully.
// - 401 Unauthorized: User authentication failed.
// - 500 Internal Server Error: Failed to retrieve expenses.
func (h *ExpenseHandler) GetExpensesList(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expenses, err := h.expenseService.GetExpensesList(r.Context(), userID)
	if err != nil {
		lib.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// GetExpensesByPeriod handles the HTTP request to retrieve expenses for the authenticated user within a specified date range.
// It expects query parameters "start" and "end" with dates in "YYYY-MM-DD" format.
// Possible HTTP responses:
// - 200 OK: Expenses retrieved successfully.
// - 400 Bad Request: Invalid or missing date parameters, or end date before start date.
// - 401 Unauthorized: User authentication failed.
// - 500 Internal Server Error: Failed to retrieve expenses.
func (h *ExpenseHandler) GetExpensesByPeriod(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, err := time.Parse("2006-01-02", r.URL.Query().Get("start"))
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid start time (use YYYY-MM-DD)")
		return
	}
	end, err := time.Parse("2006-01-02", r.URL.Query().Get("end"))
	if err != nil {
		lib.WriteJSONError(w, http.StatusBadRequest, "invalid end time (use YYYY-MM-DD)")
		return
	}

	if end.Before(start) {
		lib.WriteJSONError(w, http.StatusBadRequest, "end time must be after start time")
		return
	}

	expenses, err := h.expenseService.GetExpensesByPeriod(r.Context(), userID, start, end)
	if err != nil {
		lib.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// GetExpensesByCategory handles the HTTP request to retrieve expenses for the authenticated user filtered by category.
// It expects a query parameter "category".
// Possible HTTP responses:
// - 200 OK: Expenses retrieved successfully.
// - 400 Bad Request: Missing category parameter.
// - 401 Unauthorized: User authentication failed.
// - 500 Internal Server Error: Failed to retrieve expenses.
func (h *ExpenseHandler) GetExpensesByCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.GetUserIDFromContext(r)
	if err != nil {
		lib.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	category := r.URL.Query().Get("category")
	if category == "" {
		lib.WriteJSONError(w, http.StatusBadRequest, "category parameter is required")
		return
	}

	expenses, err := h.expenseService.GetExpensesByCategory(r.Context(), userID, category)
	if err != nil {
		lib.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}
