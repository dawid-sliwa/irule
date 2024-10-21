package middlewares

import (
	"irule-api/internal/config"
	"irule-api/internal/svc"
	"net/http"
	"strings"
)

func AuthMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.Header.Get("Authorization")
			if v == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			str := strings.Split(v, "Bearer ")
			if len(str) != 2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString := str[1]
			err := svc.VerifyToken(tokenString, cfg)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
