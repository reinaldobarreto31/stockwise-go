package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	// ContextUserID is the context key for the authenticated user's ID.
	ContextUserID contextKey = "user_id"
	// ContextUserRole is the context key for the authenticated user's role.
	ContextUserRole contextKey = "user_role"
)

// Claims holds the JWT payload.
type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth validates the Bearer token in the Authorization header.
// On success, user_id and user_role are stored in the request context.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"missing or malformed authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(os.Getenv("JWT_SECRET"))

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextUserRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext retrieves the authenticated user ID from the context.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ContextUserID).(int64)
	return id, ok
}

// UserRoleFromContext retrieves the authenticated user role from the context.
func UserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextUserRole).(string)
	return role, ok
}
