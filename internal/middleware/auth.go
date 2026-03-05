package middleware

import (
	"context"
	"expense_tracker/internal/service"
	"expense_tracker/lib"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// AuthMiddleware returns an HTTP middleware that authenticates requests using a Bearer token.
// It expects the "Authorization" header with the format "Bearer <token>".
//
// The middleware validates the token using the provided AuthService.
// If the token is valid, it extracts the user ID and stores it in the request context
// under the key defined by lib.UserIDkey, allowing subsequent handlers to identify the user.
//
// If authentication fails, it responds with HTTP 401 Unauthorized and a JSON error message.
//
// Parameters:
// - authService: a pointer to AuthService used to validate tokens.
//
// Returns:
// - A middleware function that wraps an http.Handler with authentication logic.
//
// Usage:
//
//	http.Handle("/protected", AuthMiddleware(authService)(protectedHandler))
func AuthMiddleware(authService *service.AuthService, logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.WithFields(logrus.Fields{
				"middleware": "AuthMiddleware",
				"path":       r.URL.Path,
				"method":     r.Method,
			})
			log.Debug("starting authentication")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("authorization header missing")
				lib.WriteJSONError(w, http.StatusUnauthorized, "authotization header is required")
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Warn("invalid authorization header format")
				lib.WriteJSONError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			log = log.WithField("token_length", len(tokenString))
			userID, err := authService.ValidateToken(tokenString)

			if err != nil {
				log.WithError(err).Warn("token validation failed")
				lib.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			log = log.WithField("user_id", userID)
			log.Info("user authenticated successfully")

			ctx := context.WithValue(r.Context(), lib.UserIDkey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
