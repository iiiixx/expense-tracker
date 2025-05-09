package middleware

import (
	"context"
	"expense_tracker/internal/service"
	"expense_tracker/lib"
	"net/http"
	"strings"
)

func AuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				lib.WriteJSONError(w, http.StatusUnauthorized, "authotization header is required")
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				lib.WriteJSONError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			userID, err := authService.ValidateToken(tokenString)
			if err != nil {
				lib.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), lib.UserIDkey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
