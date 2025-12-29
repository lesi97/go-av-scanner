package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/lesi97/go-av-scanner/internal/utils"
)

func ApiKey(apiKey string) func(http.Handler) http.Handler {
	expectedKey := []byte(apiKey)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := []byte(r.Header.Get("X-API-Key"))

			if len(key) != len(expectedKey) ||
				subtle.ConstantTimeCompare(key, expectedKey) != 1 {
				utils.Error(w, http.StatusUnauthorized, "invalid API key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}