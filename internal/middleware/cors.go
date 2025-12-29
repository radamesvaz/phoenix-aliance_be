package middleware

import (
	"net/http"
	"strconv"

	"phoenix-alliance-be/internal/config"
)

// CORSMiddleware adds CORS headers to responses.
// In DEV, the default is permissive (AllowAllOrigins=true when CORS_ALLOWED_ORIGINS="*").
func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(cfg.CORS.AllowedOrigins))
	for _, o := range cfg.CORS.AllowedOrigins {
		allowed[o] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Decide origin policy:
			// - If AllowCredentials=true, we cannot use "*", so we reflect the Origin (when present).
			// - If AllowAllOrigins=true and credentials are off, we can safely use "*".
			if origin != "" {
				if cfg.CORS.AllowAllOrigins {
					if cfg.CORS.AllowCredentials {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Add("Vary", "Origin")
					} else {
						w.Header().Set("Access-Control-Allow-Origin", "*")
					}
				} else {
					if _, ok := allowed[origin]; ok {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Add("Vary", "Origin")
					}
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", cfg.CORS.AllowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", cfg.CORS.AllowedHeaders)

			if cfg.CORS.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if cfg.CORS.MaxAgeSeconds > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.CORS.MaxAgeSeconds))
			}

			// Fast-path preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

