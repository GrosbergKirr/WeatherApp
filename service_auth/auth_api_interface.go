package service_auth

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

type AuthInterface interface {
	RegisterUser(log *slog.Logger, user models.User) (int, error)
	Login(log *slog.Logger, user models.User) (int, error)
}
