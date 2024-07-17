package service_auth

import (
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/go-chi/render"
)

func RegisterUser(log *slog.Logger, auth AuthInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "/service_auth/registration"
		var req models.User
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("fail to decode json. Provide login & password", slog.Any("err: ", err), slog.String("path", path))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debug("Decode json success")

		stat, err := auth.RegisterUser(log, req)
		if err != nil {
			log.Error("Fail to create user", slog.Any("err", err))
			w.WriteHeader(stat)
			return
		}
		log.Info("User created success")
		w.WriteHeader(stat)

	}
}
