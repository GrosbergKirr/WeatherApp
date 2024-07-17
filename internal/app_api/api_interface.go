package app_api

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

type DatabaseInterface interface {
	GetCities(log *slog.Logger, perPage, offset int) (models.CitiesListResponse, int, error)
	GetShortPred(log *slog.Logger, city string) (models.ShortForecastResponse, int, error)
	GetFullPred(log *slog.Logger, param models.Forecast) (models.Record, int, error)
}
