package service_auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/internal"
)

func TokenAuthMiddleware(log *slog.Logger, cfg *internal.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					log.Error("Cookies not set")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				log.Error("Get cookies error")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			key := []byte(cfg.JWTKey)
			claims, stat, err := CookieJWTGet(log, cookie, key)
			if err != nil {
				log.Error("Get data from cookies error")
				w.WriteHeader(stat)
				return
			}
			ctx := context.WithValue(r.Context(), "login", claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
			w.WriteHeader(stat)
		})
	}
}
