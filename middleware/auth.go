package middleware

import (
	"context"
	"net/http"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/model"
)

type contextKey string

const userCtxKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role") 

		user := &model.AuthUser{
			ID:   "123", 
			Role: model.Role(role),
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *model.AuthUser {
	user, _ := ctx.Value(userCtxKey).(*model.AuthUser)
	return user
}
