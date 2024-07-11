package app_api

import (
	"log/slog"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
)

type DatabaseInterface interface {
	GetCities(log *slog.Logger, perPage, offset int) (models.CitiesListResponse, error, int)
	GetShortPred(log *slog.Logger, city string) (models.ShortForecastResponse, error, int)
	GetFullPred(log *slog.Logger, param models.Forecast) (models.Record, error, int)
}
