package auth

import (
	"context"
	"net/http"

	"github.com/monachy/projek/internal/member"
)

func AuthMiddleware(authService *AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			m, err := authService.GetMemberBySession(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "member", m)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetMemberFromContext(ctx context.Context) *member.Member {
	if v := ctx.Value("member"); v != nil {
		if m, ok := v.(*member.Member); ok {
			return m
		}
	}
	return nil
}
