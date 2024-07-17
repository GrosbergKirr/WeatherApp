package service_auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
)

func LogInUser(log *slog.Logger, cfg *internal.Config, auth AuthInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "/service_auth/authentication"
		var req models.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("fail to decode json. Provide your login & password", slog.Any("err: ", err), slog.String("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debug("Decode json success")

		stat, err := auth.Login(log, req)
		if err != nil {
			log.Error("Access denied", slog.Any("err", err))
			w.WriteHeader(stat)
			return
		}
		expireAt := time.Now().Add(time.Duration(cfg.TokenExpirationTime * int(time.Minute)))
		key := []byte(cfg.JWTKey)
		claims := &Claims{
			Username: req.Login,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expireAt.Unix(),
			},
		}
		cookie, stat, err := CookieJWTCreate(log, claims, &expireAt, key)
		if err != nil {
			log.Error("Cookie jwt creation error")
			w.WriteHeader(stat)
			return
		}
		http.SetCookie(w, cookie)
		log.Info("Login to the account is successful")
		w.WriteHeader(stat)
	}
}
