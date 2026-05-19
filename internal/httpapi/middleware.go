package httpapi

import (
	"fmt"
	"net/http"
	"strings"
)

func APIKeyMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey == "" || r.URL.Path == "/healthz" || strings.HasPrefix(r.URL.Path, "/openapi") {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get("X-API-Key") == apiKey {
				next.ServeHTTP(w, r)
				return
			}
			if strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ") == apiKey {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, fmt.Sprintf("%s\n", http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		})
	}
}
