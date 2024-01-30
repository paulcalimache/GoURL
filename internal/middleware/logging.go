package middleware

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/static") {
			log.Info().Str("method", r.Method).Str("url", r.URL.Path).Send()
		}
		next.ServeHTTP(w, r)
	})
}
